package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

const (
	API_ERROR_NO_RECORD_FOUND = "no_record_found"
)

type Plugin struct {
	plugin.MattermostPlugin
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSRegion          string
	disabled           bool
}

type TranslatedMessage struct {
	Id             string `json:"id"`
	PostId         string `json:"post_id"`
	SourceLanguage string `json:"source_lang"`
	SourceText     string `json:"source_text"`
	TargetLanguage string `json:"target_lang"`
	TranslatedText string `json:"translated_text"`
	UpdateAt       int64  `json:"update_at"`
}

type UserInfo struct {
	UserID         string `json:"user_id"`
	Activated      bool   `json:"activated"`
	SourceLanguage string `json:"source_language"`
	TargetLanguage string `json:"target_language"`
}

func (p *Plugin) NewUserInfo(userId string) *UserInfo {
	return &UserInfo{
		UserID:         userId,
		Activated:      true,
		SourceLanguage: LANGUAGE_AUTO,
		TargetLanguage: LANGUAGE_EN,
	}
}

func (u *UserInfo) IsValid() error {
	if u.UserID == "" || len(u.UserID) != 26 {
		return fmt.Errorf("Invalid: user_id field")
	}

	if u.SourceLanguage == "" {
		return fmt.Errorf("Invalid: source_language field")
	}

	if u.TargetLanguage == "" {
		return fmt.Errorf("Invalid: target_language field")
	}

	if LANGUAGE_CODES[u.SourceLanguage] == nil {
		return fmt.Errorf("Invalid: source_language must be in a supported language code")
	}

	if LANGUAGE_CODES[u.TargetLanguage] == nil {
		return fmt.Errorf("Invalid: target_language must be in a supported language code")
	}

	if u.SourceLanguage == u.TargetLanguage {
		return fmt.Errorf("Invalid: source_language and target_language are equal")
	}

	if u.SourceLanguage == LANGUAGE_EN && u.TargetLanguage == LANGUAGE_EN {
		return fmt.Errorf("Invalid: source_language and target_language should not be both English")
	}

	if u.SourceLanguage != LANGUAGE_EN && u.TargetLanguage != LANGUAGE_EN {
		return fmt.Errorf("Invalid: Either source_language or target_language should be English")
	}

	if u.SourceLanguage == LANGUAGE_AUTO && u.TargetLanguage != LANGUAGE_EN {
		return fmt.Errorf("Invalid: if source_language is auto, target_language must be en")
	}

	if u.TargetLanguage == LANGUAGE_AUTO {
		return fmt.Errorf("Invalid: target_language must not be auto")
	}

	return nil
}

func (p *Plugin) OnActivate() error {
	if err := p.IsValid(); err != nil {
		return err
	}

	p.API.RegisterCommand(getCommand())

	return nil
}

func (p *Plugin) IsValid() error {
	if p.AWSAccessKeyID == "" {
		return fmt.Errorf("Must have AWS Access Key ID")
	}

	if p.AWSSecretAccessKey == "" {
		return fmt.Errorf("Must have AWS Secret Access Key")
	}

	if p.AWSRegion == "" {
		return fmt.Errorf("AWS Region")
	}

	return nil
}

func (p *Plugin) getUserInfo(userID string) (*UserInfo, *APIErrorResponse) {
	var userInfo UserInfo

	if infoBytes, err := p.API.KVGet(userID); err != nil || infoBytes == nil {
		return nil, &APIErrorResponse{ID: API_ERROR_NO_RECORD_FOUND, Message: "No record found.", StatusCode: http.StatusBadRequest}
	} else if err := json.Unmarshal(infoBytes, &userInfo); err != nil {
		return nil, &APIErrorResponse{ID: "unable_to_unmarshal", Message: "Unable to unmarshal json.", StatusCode: http.StatusInternalServerError}
	}

	return &userInfo, nil
}

func (p *Plugin) setUserInfo(userInfo *UserInfo) *APIErrorResponse {
	if err := userInfo.IsValid(); err != nil {
		return &APIErrorResponse{ID: "invalid_user_info", Message: err.Error(), StatusCode: http.StatusBadRequest}
	}

	jsonUserInfo, err := json.Marshal(userInfo)
	if err != nil {
		return &APIErrorResponse{ID: "unable_to_unmarshal", Message: "Unable to marshal json.", StatusCode: http.StatusInternalServerError}
	}

	if err := p.API.KVSet(userInfo.UserID, jsonUserInfo); err != nil {
		return &APIErrorResponse{ID: "unable_to_save", Message: "Unable to save user info.", StatusCode: http.StatusInternalServerError}
	}

	p.emitUserInfoChange(userInfo)

	return nil
}

func (u *UserInfo) getActivatedString() string {
	activated := "off"
	if u.Activated {
		activated = "on"
	}

	return activated
}

func (p *Plugin) emitUserInfoChange(userInfo *UserInfo) {
	p.API.PublishWebSocketEvent(
		"info_change",
		map[string]interface{}{
			"user_id":         userInfo.UserID,
			"activated":       userInfo.Activated,
			"source_language": userInfo.SourceLanguage,
			"target_language": userInfo.TargetLanguage,
		},
		&model.WebsocketBroadcast{UserId: userInfo.UserID},
	)
}
