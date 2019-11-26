package learn

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// LearnResponse is a simple struct defining the shape of data we care about
// that comes back from notifying Learn for decoding into.
type LearnResponse struct {
	ReleaseID         int               `json:"release_id"`
	PreviewURL        string            `json:"preview_url"`
	Errors            string            `json:"errors"`
	Status            string            `json:"status"`
	GLearnCredentials GLearnCredentials `json:"glearn_credentials"`
}

// GLearnCredentials represents the important AWS credentials we retrieve from Learn
// with an api_token
type GLearnCredentials struct {
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	KeyPrefix       string `json:"key_prefix"`
	BucketName      string `json:"bucket_name"`
}

func (api *ApiClient) PollForBuildResponse(releaseID int, attempts *uint8) (*LearnResponse, error) {
	client := &http.Client{Timeout: time.Second * 30}

	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:3003/api/v1/releases/%d/release_polling", releaseID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.token))

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error: response status: %d", res.StatusCode)
	}

	var l LearnResponse
	err = json.NewDecoder(res.Body).Decode(&l)
	if err != nil {
		return nil, err
	}

	if l.Status == "processing" || l.Status == "pending" {
		*attempts--
		time.Sleep(2 * time.Second)

		if *attempts == uint8(0) {
			return nil, errors.New(
				"Sorry, we are having trouble requesting your preview build from Learn. Please try again",
			)
		}

		return pollForBuildResponse(releaseID, attempts)
	}

	return &l, nil
}

// notifyLearn takes an s3 bucket key name as an argument is used to tell Learn there is new preview
// content on s3 and where to find it so it can build/preview.
func (api *ApiClient) NotifyLearn(bucketKey string, isDirectory bool) (*LearnResponse, error) {
	payload := map[string]string{
		"s3_key": bucketKey,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	var endpoint string
	if isDirectory {
		endpoint = "/api/v1/releases"
	} else {
		endpoint = "/api/v1/content_file"
	}

	client := &http.Client{Timeout: time.Second * 30}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("http://localhost:3003%s", endpoint),
		bytes.NewBuffer(payloadBytes),
	)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.token))

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error: response status: %d", res.StatusCode)
	}

	l := &LearnResponse{}
	json.NewDecoder(res.Body).Decode(l)

	return l, nil
}

// retrieveS3CredentialsWithAPIKey uses a user's api_token to request AWS credentials
// from Learn. It returns a populated *GLearnCredentials struct or an error
func (api *ApiClient) RetrieveS3CredentialsWithAPIKey() (*GLearnCredentials, error) {
	apiToken, ok := viper.Get("api_token").(string)
	if !ok {
		return nil, errors.New("Please set your api_token in ~/.glearn-config.yaml")
	}

	client := &http.Client{Timeout: time.Second * 30}

	req, err := http.NewRequest("GET", "http://localhost:3003/api/v1/users/glearn_credentials", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.token))

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var l LearnResponse
	err = json.NewDecoder(res.Body).Decode(&l)
	if err != nil {
		return nil, err
	}

	return &GLearnCredentials{
		AccessKeyID:     l.GLearnCredentials.AccessKeyID,
		SecretAccessKey: l.GLearnCredentials.SecretAccessKey,
		KeyPrefix:       l.GLearnCredentials.KeyPrefix,
		BucketName:      l.GLearnCredentials.BucketName,
	}, nil
}
