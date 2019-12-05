package learn

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gSchool/glearn-cli/api"
	"github.com/spf13/viper"
)

// API is the exported APIClient, it is set during Init
var API *APIClient

// APIClient makes network API calls to Learn
type APIClient struct {
	client      api.Client
	baseURL     string
	Credentials *Credentials
}

// Credentials represents the shape of data that the initial call to Learn
// for s3 and slack credentials will hydrate
type Credentials struct {
	*APIToken         `json:"api_token"`
	*S3Credentials    `json:"s3_credentials"`
	*SlackCredentials `json:"slack_credentials"`
}

// S3Credentials represents the important AWS credentials we retrieve from Learn
// with an api_token
type S3Credentials struct {
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	KeyPrefix       string `json:"key_prefix"`
	BucketName      string `json:"bucket_name"`
}

// SlackCredentials represents the credentials we retrieve from Learn for the CLI
// to operate correctly
type SlackCredentials struct {
	DevNotifyUrl string `json:"dev_notify_url"`
}

// APIToken is a simple wrapper around an API token
type APIToken struct {
	token string
}

// CredentialsResponse describes the shape of the return data from the call
// to RetrieveCredentials
type CredentialsResponse struct {
	S3    S3Credentials    `json:"s3"`
	Slack SlackCredentials `json:"slack"`
}

// NewAPI is a constructor for the ApiClient
func NewAPI(baseURL string, client api.Client) (*APIClient, error) {
	// Retrieve the application credentials for the CLI using a user's API token
	creds, err := API.RetrieveCredentials()
	if err != nil {
		return nil, errors.New("Could not retrieve credentials from Learn. Please ensure you have the right API token in your ~/.glearn-config.yaml")
	}

	return &APIClient{
		client:      client,
		baseURL:     baseURL,
		Credentials: creds,
	}, nil
}

// RetrieveCredentials uses a user's api_token to request AWS credentials
// from Learn. It returns a populated *S3Credentials struct or an error
func (api *APIClient) RetrieveCredentials() (*Credentials, error) {
	// Early return if user's api_token is not set
	apiToken, ok := viper.Get("api_token").(string)
	if !ok {
		return nil, errors.New("Please set your api_token in ~/.glearn-config.yaml")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/users/learn_cli_credentials", api.baseURL), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api.Credentials.token))

	res, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var c CredentialsResponse

	err = json.NewDecoder(res.Body).Decode(&c)
	if err != nil {
		return nil, err
	}

	return &Credentials{
		S3Credentials: &S3Credentials{
			AccessKeyID:     c.S3.AccessKeyID,
			SecretAccessKey: c.S3.SecretAccessKey,
			KeyPrefix:       c.S3.KeyPrefix,
			BucketName:      c.S3.BucketName,
		},
		SlackCredentials: &SlackCredentials{
			DevNotifyUrl: c.Slack.DevNotifyUrl,
		},
		APIToken: &APIToken{apiToken},
	}, nil
}

// NotifySlack is used throughout the CLI for production error handling
func (api *APIClient) NotifySlack(err error) {
	// Do not notify slack during development
	if api.Credentials.DevNotifyUrl == "development" {
		return
	}

	go func(err error) {
		textMsg := struct {
			Text string `json:"text"`
		}{
			Text: fmt.Sprintf("%s: %s", api.baseURL, err),
		}

		bytePostData, err := json.Marshal(textMsg)

		req, err := http.NewRequest("POST", api.Credentials.DevNotifyUrl, bytes.NewReader(bytePostData))
		if err == nil {
			req.Header.Add("Content-Type", "application/json; charset=utf-8")
			client := &http.Client{Timeout: time.Second * 30}
			client.Do(req)
		}
	}(err)
}
