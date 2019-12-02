package learn

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// PreviewResponse is a simple struct defining the shape of data we care about
// that comes back from notifying Learn for decoding into.
type PreviewResponse struct {
	ReleaseID      int            `json:"release_id"`
	PreviewURL     string         `json:"preview_url"`
	Errors         string         `json:"errors"`
	Status         string         `json:"status"`
	S3Crednentials S3Crednentials `json:"glearn_credentials"`
}

// S3Crednentials represents the important AWS credentials we retrieve from Learn
// with an api_token
type S3Crednentials struct {
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
	KeyPrefix       string `json:"key_prefix"`
	BucketName      string `json:"bucket_name"`
}

// PollForBuildResponse attempts to check if a release has finished building every 2 seconds.
func (api *APIClient) PollForBuildResponse(releaseID int, attempts *uint8) (*PreviewResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/releases/%d/release_polling", api.baseURL, releaseID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api.token))

	res, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var p PreviewResponse

	if res.StatusCode != http.StatusOK {
		json.NewDecoder(res.Body).Decode(&p)
		return nil, fmt.Errorf("Error: response status: %d, response body: %s", res.StatusCode, p.Errors)
	}

	err = json.NewDecoder(res.Body).Decode(&p)
	if err != nil {
		return nil, err
	}

	if p.Status == "processing" || p.Status == "pending" {
		*attempts--
		time.Sleep(2 * time.Second)

		if *attempts == uint8(0) {
			return nil, errors.New(
				"Sorry, we are having trouble requesting your build from Learn. Please try again",
			)
		}

		return api.PollForBuildResponse(releaseID, attempts)
	}

	return &p, nil
}

// BuildReleaseFromS3 takes an s3 bucket key name as an argument is used to tell Learn there is new preview
// content on s3 and where to find it so it can build/preview.
func (api *APIClient) BuildReleaseFromS3(bucketKey string, isDirectory bool) (*PreviewResponse, error) {
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
		endpoint = "/api/v1/content_files"
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s%s", api.baseURL, endpoint),
		bytes.NewBuffer(payloadBytes),
	)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api.token))

	res, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error: response status: %d", res.StatusCode)
	}

	p := &PreviewResponse{}
	json.NewDecoder(res.Body).Decode(p)

	return p, nil
}

// RetrieveS3Credentials uses a user's api_token to request AWS credentials
// from Learn. It returns a populated *S3Crednentials struct or an error
func (api *APIClient) RetrieveS3Credentials() (*S3Crednentials, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/users/glearn_credentials", api.baseURL), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api.token))

	res, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var p PreviewResponse
	err = json.NewDecoder(res.Body).Decode(&p)
	if err != nil {
		return nil, err
	}

	return &S3Crednentials{
		AccessKeyID:     p.S3Crednentials.AccessKeyID,
		SecretAccessKey: p.S3Crednentials.SecretAccessKey,
		KeyPrefix:       p.S3Crednentials.KeyPrefix,
		BucketName:      p.S3Crednentials.BucketName,
	}, nil
}
