package main

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

const (
	TRANSLATE_USERNAME = "Autotranslate Plugin"
	TRANSLATE_ICON_URL = "https://docs.mattermost.com/_images/icon-76x76.png"
	LANGUAGE_EN        = "en"
	LANGUAGE_AUTO      = "auto"
)

const COMMAND_HELP = `
This autotranslation plugin is powered by Amazon Translate which supports translation between English and any of the following languages: Arabic, Chinese (Simplified), Chinese (Traditional), Czech, French, German, Italian, Japanese, Portuguese, Russian, Spanish, and Turkish.

* |/autotranslate on| - Add an option to translate a post and will have default setting of Auto as source and English as target.
* |/autotranslate off| - Remove an option to translate a post
* |/autotranslate info| - Show user info on this plugin
* |/autotranslate source [value]| - Update your autotranslation source
  * |value| can be any of the supported language codes below.
  * Note: Changing source setting will automatically update target into English.
     * Ex. |/autotranslate source ar| will set source to Arabic and target to English
    * Ex. |/autotranslate source auto| will set source to Auto and target to English
* |/autotranslate target [value]| - Update your autotranslation target
  * |value| can be any of the supported language codes below except auto.
  * Note: In most cases, changing target setting will automatically update source into English.
    * Ex. |/autotranslate target ar| will set target to Arabic and source to English
    * Ex. |/autotranslate target auto| will not change target settings and will return an error
    * Ex. |/autotranslate target en| will set target to English and automatically set source to Auto
* |Language codes|:
  * auto (Auto) : Automatic detection based on supported language below
  * ar (Arabic)
  * zh (Chinese)
  * cs (Czech)
  * fr (French)
  * de (German)
  * en (English)
  * es (Spanish)
  * it (Italian)
  * ja (Japanese)
  * pt (Portuguese)
  * ru (Russian)
  * tr (Turkish)
  `

var LANGUAGE_CODES = map[string]interface{}{
	"auto": "Auto",
	"ar":   "Arabic",
	"zh":   "Chinese",
	"cs":   "Czech",
	"fr":   "French",
	"de":   "German",
	"en":   "English",
	"es":   "Spanish",
	"it":   "Italian",
	"ja":   "Japanese",
	"pt":   "Portuguese",
	"ru":   "Russian",
	"tr":   "Turkish",
}

func (p *Plugin) registerCommands() error {
	if err := p.API.RegisterCommand(&model.Command{
		Trigger:          "autotranslate",
		DisplayName:      "Autotranslate",
		Description:      "Mattermost Autotranslation Plugin",
		AutoComplete:     true,
		AutoCompleteDesc: "Available commands: info, on, off, source, target, help",
		AutoCompleteHint: "[command]",
	}); err != nil {
		return errors.Wrap(err, "failed to register autotranslate command")
	}

	return nil
}

func getCommandResponse(responseType, text string) *model.CommandResponse {
	return &model.CommandResponse{
		ResponseType: responseType,
		Text:         text,
		Username:     TRANSLATE_USERNAME,
		IconURL:      TRANSLATE_ICON_URL,
		Type:         model.POST_DEFAULT,
	}
}

func setUserInfoCommandResponse(userInfo *UserInfo, err *APIErrorResponse, action string) (*model.CommandResponse, *model.AppError) {
	var actionMapping = map[string]interface{}{
		"source": "setting up language source of autotranslation plugin",
		"target": "setting up language target of autotranslation plugin",
		"on":     "turning on the autotranslation plugin",
		"off":    "turning off the autotranslation plugin",
		"info":   "getting user information",
	}

	if err != nil {
		errorMessage := ""
		if len(err.Message) > 0 {
			errorMessage = err.Message
		}

		text := fmt.Sprintf("An error occurred %s. `%s`", actionMapping[action], errorMessage)
		if err.ID == API_ERROR_NO_RECORD_FOUND {
			text = "No record found. If not yet turned on for the first time, try `/autotranslate on` to enable."
		}

		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, text), nil
	}

	text := fmt.Sprintf(
		"Successfully updated!\nYour autotranslation plugin settings:\n * Active: `%s`\n * Language: `source: %s`, `target: %s`\n",
		userInfo.getActivatedString(), LANGUAGE_CODES[userInfo.SourceLanguage], LANGUAGE_CODES[userInfo.TargetLanguage],
	)

	if action == "off" {
		text = "Autotranslate plugin is turned off."
	}

	return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, text), nil
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	split := strings.Fields(args.Command)
	command := split[0]
	action := ""
	param := ""
	if len(split) > 1 {
		action = split[1]
	}
	if len(split) > 2 {
		param = split[2]
	}

	if command != "/autotranslate" {
		return nil, nil
	}

	var text = ""
	if action == "" || action == "help" {
		text = "###### Mattermost Autotranslate Plugin - Slash Command Help\n" + strings.Replace(COMMAND_HELP, "|", "`", -1)
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, text), nil
	}

	userInfo, err := p.getUserInfo(args.UserId)
	if userInfo == nil && action != "on" {
		text = "No record found. Try `/autotranslate on` to enable."
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, text), nil
	}

	switch action {
	case "info":
		return setUserInfoCommandResponse(userInfo, err, action)
	case "on":
		if userInfo == nil {
			userInfo = p.NewUserInfo(args.UserId)
		} else {
			userInfo.Activated = true
		}

		err = p.setUserInfo(userInfo)
		return setUserInfoCommandResponse(userInfo, err, action)
	case "off":
		if userInfo == nil {
			return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "No record found. If not yet turned on for the first time, try `/autotranslate on` to enable. Otherwise, your record is lost for unknown reason."), nil
		}

		userInfo.Activated = false
		err = p.setUserInfo(userInfo)
		return setUserInfoCommandResponse(userInfo, err, action)
	case "source":
		if userInfo == nil {
			return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "No record found. If not yet turned on for the first time, try `/autotranslate on` to enable. Otherwise, your record is lost for unknown reason."), nil
		}

		if param == "" || LANGUAGE_CODES[param] == nil {
			return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Invalid parameter. Shoud pass a valid language code to change the source language"), nil
		}

		if param == userInfo.TargetLanguage { // switch language setting
			oldSource := userInfo.SourceLanguage
			userInfo.SourceLanguage = param
			userInfo.TargetLanguage = oldSource
		} else if param != LANGUAGE_EN {
			userInfo.SourceLanguage = param
			userInfo.TargetLanguage = LANGUAGE_EN
		}

		err = p.setUserInfo(userInfo)
		return setUserInfoCommandResponse(userInfo, err, action)
	case "target":
		if userInfo == nil {
			return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "No record found. If not yet turned on for the first time, try `/autotranslate on` to enable."), nil
		}

		if param == "" || LANGUAGE_CODES[param] == nil {
			return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, "Invalid parameter. Shoud pass a valid language code to change the target language"), nil
		}

		if param != LANGUAGE_AUTO && param == userInfo.SourceLanguage { // switch language setting
			oldTarget := userInfo.TargetLanguage
			userInfo.TargetLanguage = param
			userInfo.SourceLanguage = oldTarget
		} else if param == LANGUAGE_EN {
			userInfo.TargetLanguage = LANGUAGE_EN
			userInfo.SourceLanguage = LANGUAGE_AUTO
		} else {
			userInfo.TargetLanguage = param
			userInfo.SourceLanguage = LANGUAGE_EN
		}

		err = p.setUserInfo(userInfo)
		return setUserInfoCommandResponse(userInfo, err, action)
	default:
		text = "###### Mattermost Autotranslate Plugin - Slash Command Help\n" + strings.Replace(COMMAND_HELP, "|", "`", -1)
		return getCommandResponse(model.COMMAND_RESPONSE_TYPE_EPHEMERAL, text), nil
	}
}
