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
var LearnUserId string
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
	DevNotifyURL string `json:"dev_notify_url"`
}

// APIToken is a simple wrapper around an API token
type APIToken struct {
	token string
}

// CredentialsResponse describes the shape of the return data from the call
// to RetrieveCredentials
type CredentialsResponse struct {
	UserId string           `json:"user_id"`
	Email  string           `json:"user_email"`
	S3     S3Credentials    `json:"s3"`
	Slack  SlackCredentials `json:"slack"`
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
func NewAPI(baseURL string, client api.Client) (*APIClient, error) {
	apiClient := &APIClient{
		client:  client,
		baseURL: baseURL,
	}

	// Retrieve the application credentials for the CLI using a user's API token
	creds, err := apiClient.RetrieveCredentials()
	if err != nil {
		return nil, errors.New(
			fmt.Sprintf("Could not retrieve credentials from Learn. Please reset your API token with this command: learn set --api_token=your-token-from-%s/api_token", baseURL),
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
func (api *APIClient) RetrieveCredentials() (*Credentials, error) {
	// Early return if user's api_token is not set
	apiToken, ok := viper.Get("api_token").(string)
	if !ok {
		return nil, errors.New("Please set your API token with this command: learn set --api_token=your-token-from-https://learn-2.galvanize.com/api_token")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/users/learn_cli_credentials", api.baseURL), nil)
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

	LearnUserId = c.UserId
	LearnUserEmail = c.Email

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
			DevNotifyURL: c.Slack.DevNotifyURL,
		},
		APIToken: &APIToken{apiToken},
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
		Text: fmt.Sprintf("UserId: %s\nUserEmail: %s\n%s", LearnUserId, LearnUserEmail, err),
	}

	bytePostData, _ := json.Marshal(msg)

	req, err := http.NewRequest("POST", api.Credentials.DevNotifyURL, bytes.NewReader(bytePostData))
	if err == nil {
		req.Header.Add("Content-Type", "application/json; charset=utf-8")
		client := &http.Client{Timeout: time.Second * 30}
		client.Do(req)
	}
}
