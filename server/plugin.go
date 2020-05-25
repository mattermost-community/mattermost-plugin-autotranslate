package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

const (
	apiErrorNoRecordFound = "no_record_found"
)

// Plugin is a collection of fields for plugin
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

// TranslatedMessage is a collection of fields for translated message
type TranslatedMessage struct {
	ID             string `json:"id"`
	PostID         string `json:"post_id"`
	SourceLanguage string `json:"source_lang"`
	SourceText     string `json:"source_text"`
	TargetLanguage string `json:"target_lang"`
	TranslatedText string `json:"translated_text"`
	UpdateAt       int64  `json:"update_at"`
}

// UserInfo is a collection of fields for user info
type UserInfo struct {
	UserID         string `json:"user_id"`
	Activated      bool   `json:"activated"`
	SourceLanguage string `json:"source_language"`
	TargetLanguage string `json:"target_language"`
}

// NewUserInfo returns new user info
func (p *Plugin) NewUserInfo(userID string) *UserInfo {
	return &UserInfo{
		UserID:         userID,
		Activated:      true,
		SourceLanguage: autoLanguage,
		TargetLanguage: enLanguage,
	}
}

// IsValid validates user information
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

	if languageCodes[u.SourceLanguage] == "" {
		return fmt.Errorf("Invalid: source_language must be in a supported language code")
	}

	if languageCodes[u.TargetLanguage] == "" {
		return fmt.Errorf("Invalid: target_language must be in a supported language code")
	}

	if u.SourceLanguage == u.TargetLanguage {
		return fmt.Errorf("Invalid: source_language and target_language are equal")
	}

	if u.TargetLanguage == autoLanguage {
		return fmt.Errorf("Invalid: target_language must not be \"auto\"")
	}

	return nil
}

func (p *Plugin) getUserInfo(userID string) (*UserInfo, *APIErrorResponse) {
	var userInfo UserInfo

	if infoBytes, err := p.API.KVGet(userID); err != nil || infoBytes == nil {
		return nil, &APIErrorResponse{ID: apiErrorNoRecordFound, Message: "No record found.", StatusCode: http.StatusBadRequest}
	} else if err := json.Unmarshal(infoBytes, &userInfo); err != nil {
		return nil, &APIErrorResponse{ID: "unable_to_unmarshal", Message: "Unable to unmarshal json.", StatusCode: http.StatusBadRequest}
	}

	return &userInfo, nil
}

func (p *Plugin) setUserInfo(userInfo *UserInfo) *APIErrorResponse {
	if err := userInfo.IsValid(); err != nil {
		return &APIErrorResponse{ID: "invalid_user_info", Message: err.Error(), StatusCode: http.StatusBadRequest}
	}

	jsonUserInfo, err := json.Marshal(userInfo)
	if err != nil {
		return &APIErrorResponse{ID: "unable_to_unmarshal", Message: "Unable to marshal json.", StatusCode: http.StatusBadRequest}
	}

	if err := p.API.KVSet(userInfo.UserID, jsonUserInfo); err != nil {
		return &APIErrorResponse{ID: "unable_to_save", Message: "Unable to save user info.", StatusCode: http.StatusBadRequest}
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
