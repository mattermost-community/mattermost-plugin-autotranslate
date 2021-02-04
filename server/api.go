package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mattermost/mattermost-server/v5/plugin"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/translate"
)

// APIErrorResponse as standard response error
type APIErrorResponse struct {
	ID         string `json:"id"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func writeAPIError(w http.ResponseWriter, err *APIErrorResponse) {
	b, _ := json.Marshal(err)
	w.WriteHeader(err.StatusCode)
	w.Write(b)
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	if err := p.IsValid(); err != nil {
		http.Error(w, "This plugin is not configured.", http.StatusNotImplemented)
	}

	w.Header().Set("Content-Type", "application/json")

	switch path := r.URL.Path; path {
	case "/api/go":
		p.getGo(w, r)
	case "/api/get_info":
		p.getInfo(w, r)
	case "/api/set_info":
		p.setInfo(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (p *Plugin) getGo(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		http.Error(w, "Not authorized to translate post", http.StatusUnauthorized)
		return
	}

	postID := r.URL.Query().Get("post_id")
	if len(postID) != 26 {
		http.Error(w, "Invalid parameter: post_id", http.StatusBadRequest)
		return
	}

	source := r.URL.Query().Get("source")
	if len(source) < 2 || len(source) > 4 {
		http.Error(w, "Invalid parameter: source", http.StatusBadRequest)
		return
	}

	target := r.URL.Query().Get("target")
	if len(target) < 2 || len(target) > 5 {
		http.Error(w, "Invalid parameter: target", http.StatusBadRequest)
		return
	}

	post, err := p.API.GetPost(postID)
	if err != nil {
		http.Error(w, "No post to translate", http.StatusBadRequest)
		return
	}

	configuration := p.getConfiguration()
	sess := session.Must(session.NewSession())
	creds := credentials.NewStaticCredentials(configuration.AWSAccessKeyID, configuration.AWSSecretAccessKey, "")
	_, awsErr := creds.Get()
	if awsErr != nil {
		http.Error(w, "Bad credentials", http.StatusForbidden)
		return
	}

	svc := translate.New(sess, aws.NewConfig().WithCredentials(creds).WithRegion(configuration.AWSRegion))

	input := translate.TextInput{
		SourceLanguageCode: &source,
		TargetLanguageCode: &target,
		Text:               &post.Message,
	}

	output, awsErr := svc.Text(&input)
	if awsErr != nil {
		http.Error(w, awsErr.Error(), http.StatusBadRequest)
		return
	}

	translated := TranslatedMessage{
		ID:             postID + source + target + strconv.FormatInt(post.UpdateAt, 10),
		PostID:         postID,
		SourceLanguage: source,
		SourceText:     post.Message,
		TargetLanguage: target,
		TranslatedText: *output.TranslatedText,
		UpdateAt:       post.UpdateAt,
	}

	resp, _ := json.Marshal(translated)
	w.Write(resp)
}

func (p *Plugin) getInfo(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		// silently return as user is probably not logged in
		return
	}

	info, err := p.getUserInfo(userID)
	if err != nil {
		// silently return as user may not have activated the autotranslation
		return
	}

	resp, _ := json.Marshal(info)
	w.Write(resp)
}

func (p *Plugin) setInfo(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		http.Error(w, "Not authorized to set info", http.StatusUnauthorized)
		return
	}

	var info *UserInfo
	json.NewDecoder(r.Body).Decode(&info)
	if info == nil {
		http.Error(w, "Invalid parameter: info", http.StatusBadRequest)
		return
	}

	if err := info.IsValid(); err != nil {
		http.Error(w, fmt.Sprintf("Invalid info: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if info.UserID != userID {
		http.Error(w, "Invalid parameter: user mismatch", http.StatusBadRequest)
		return
	}

	err := p.setUserInfo(info)
	if err != nil {
		http.Error(w, "Failed to set info", http.StatusBadRequest)
		return
	}

	resp, _ := json.Marshal(info)
	w.Write(resp)
}
