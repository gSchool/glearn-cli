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

// Learn data used for reporting
var LearnUserId int
var LearnUserEmail string

// APIClient makes network API calls to Learn
type APIClient struct {
	client      api.Client
	baseURL     string
	Credentials *Credentials
}

// Credentials represents the shape of data that the initial call to Learn
// for s3 and slack credentials will hydrate
type Credentials struct {
	*APIToken    `json:"api_token"`
	DevNotifyURL string `json:"dev_notify_url"`
	PresignedUrl string `json:"presigned_url"`
	S3Key        string `json:"s3_key"`
	UserId       int    `json:"user_id"`
}

// APIToken is a simple wrapper around an API token
type APIToken struct {
	token string
}

// CredentialsResponse describes the shape of the return data from the call
// to RetrieveCredentials
type CredentialsResponse struct {
	UserId       int    `json:"user_id"`
	Email        string `json:"user_email"`
	PresignedUrl string `json:"presigned_url"`
	S3Key        string `json:"s3_key"`
	DevNotifyUrl string `json:"dev_notify_url"`
}

// CLIBenchmarkPayload is the shape of the payload to send to Learn's learn_cli_metadata
// endpoint.
type CLIBenchmarkPayload struct {
	*CLIBenchmark `json:"cli_benchmark"`
}

// CLIBenchmark holds timing in ms for the 3 main actions in the preview command
type CLIBenchmark struct {
	// All in millisconds
	Compression           int64  `json:"time_to_compress,omitempty"`
	UploadToS3            int64  `json:"time_to_upload_to_s3,omitempty"`
	LearnBuild            int64  `json:"time_to_build_on_learn,omitempty"`
	MasterReleaseAndBuild int64  `json:"master_release_and_build,omitempty"`
	TotalCmdTime          int64  `json:"total_cmd_time,omitempty"`
	CmdName               string `json:"command_name,omitempty"`
}

// NewAPI is a constructor for the ApiClient
func NewAPI(baseURL string, client api.Client, getPresignedPostUrl bool) (*APIClient, error) {
	apiClient := &APIClient{
		client:  client,
		baseURL: baseURL,
	}

	// Retrieve the application credentials for the CLI using a user's API token
	creds, err := apiClient.RetrieveCredentials(getPresignedPostUrl)
	if err != nil {
		return nil, fmt.Errorf(
			"Could not retrieve credentials from Learn. Please reset your API token with this command: learn set --api_token=your-token-from-%s/api_token\n\n", baseURL,
		)
	}

	apiClient.Credentials = creds

	return apiClient, nil
}

// BaseURL returns the clients baseURL
func (api *APIClient) BaseURL() string {
	return api.baseURL
}

// RetrieveCredentials uses a user's api_token to request AWS credentials
// from Learn. It returns a populated *S3Credentials struct or an error
func (api *APIClient) RetrieveCredentials(getPresignedPostUrl bool) (*Credentials, error) {
	// Early return if user's api_token is not set
	apiToken, ok := viper.Get("api_token").(string)
	if !ok {
		return nil, errors.New("Please set your API token with this command: learn set --api_token=your-token-from-https://learn-2.galvanize.com/api_token")
	}

	// add presignedParam to request one in the response
	presignedParam := ""
	if getPresignedPostUrl {
		presignedParam = "?presigned_url=true"
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/users/cli_access%s", api.baseURL, presignedParam), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Source", "gLearn_cli")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))

	res, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error: response status: %d", res.StatusCode)
	}

	var c CredentialsResponse

	err = json.NewDecoder(res.Body).Decode(&c)
	if err != nil {
		return nil, err
	}

	LearnUserId = c.UserId
	LearnUserEmail = c.Email

	return &Credentials{
		DevNotifyURL: c.DevNotifyUrl,
		PresignedUrl: c.PresignedUrl,
		S3Key:        c.S3Key,
		APIToken:     &APIToken{apiToken},
		UserId:       c.UserId,
	}, nil
}

// SendMetadataToLearn takes a *CLIBenchmarkPayload struct payload to send to Learn
// for monitoring how long everything is taking
func (api *APIClient) SendMetadataToLearn(timingPayload *CLIBenchmarkPayload) error {
	payloadBytes, err := json.Marshal(timingPayload)
	if err != nil {
		return err
	}

	endpoint := "/api/v1/users/learn_cli_metadata"

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s%s", api.baseURL, endpoint),
		bytes.NewBuffer(payloadBytes),
	)
	if err != nil {
		return err
	}
	defer req.Body.Close()

	req.Header.Set("Source", "gLearn_cli")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api.Credentials.token))

	res, err := api.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Error: response status: %d", res.StatusCode)
	}

	return nil
}

// NotifySlack is used throughout the CLI for production error handling
func (api *APIClient) NotifySlack(err error) {
	// Do not notify slack during development
	if api.Credentials.DevNotifyURL == "development" {
		return
	}

	msg := struct {
		Text string `json:"text"`
	}{
		Text: fmt.Sprintf("UserId: %d\nUserEmail: %s\n%s", LearnUserId, LearnUserEmail, err),
	}

	bytePostData, _ := json.Marshal(msg)

	req, err := http.NewRequest("POST", api.Credentials.DevNotifyURL, bytes.NewReader(bytePostData))
	if err == nil {
		req.Header.Add("Content-Type", "application/json; charset=utf-8")
		client := &http.Client{Timeout: time.Second * 30}
		client.Do(req)
	}
}
