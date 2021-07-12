package learn

import (
	"fmt"
	"testing"

	"github.com/gSchool/glearn-cli/api"
	"github.com/spf13/viper"
)

const validPreviewResponse = `{"status":"success","release_id":1,"preview_url":"http://example.com"}`
const pendingPreviewResponse = `{"status":"pending","release_id":1,"preview_url":"http://example.com"}`
const credentialsResponse = `{"s3":{"access_key_id":"access_keyin","secret_access_key":"secret_keyin","key_prefix":"keykey's delivery service","bucket_name":"buqet"}, "slack":{"dev_notify_url": "development"}}`

func Test_PollForBuildResponse(t *testing.T) {
	viper.Set("api_token", "apiToken")
	mockClient := api.MockResponse(validPreviewResponse)
	API, _ := NewAPI("https://example.com", mockClient)

	attempts := uint8(1)
	previewResponse, err := API.PollForBuildResponse(1, false, "foo.md", &attempts)
	if err != nil {
		t.Errorf("error not nil: %s\n", err)
	}
	if previewResponse.ReleaseID != 1 {
		t.Errorf("Failed to properly json parse the preview response body")
	}

	// verify that requests were made properly
	if len(mockClient.Requests) != 2 {
		t.Errorf("fetching the block should make two requests")
		return
	}

	req := mockClient.Requests[0]
	if req.Method != "GET" {
		t.Errorf("Request made to Learn should be a GET, was %s", req.Method)
	}

	req = mockClient.Requests[1]
	if req.Method != "GET" {
		t.Errorf("Request made to Learn should be a GET, was %s", req.Method)
	}

	urlTarget := "https://example.com/api/v1/releases/1/release_polling?context=foo.md"
	if req.URL.String() != urlTarget {
		t.Errorf("Request made to Learn should be to url '%s' but was '%s'\n", urlTarget, req.URL.String())
	}
	if req.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type header should be 'application/json', was '%s'\n", req.Header.Get("Content-Type"))
	}
	if req.Header.Get("Authorization") != "Bearer apiToken" {
		t.Errorf("Authorization header should be 'Basic apiToken', was '%s'\n", req.Header.Get("Authorization"))
	}
}

func Test_PollForBuildResponse_EndAttempts(t *testing.T) {
	viper.Set("api_token", "apiToken")
	mockClient := api.MockResponse(pendingPreviewResponse)
	API, _ := NewAPI("https://example.com", mockClient)

	attempts := uint8(1)
	_, err := API.PollForBuildResponse(1, true, "", &attempts)
	if err == nil {
		t.Errorf("error should be present if attempts are exausted nil: %s\n", err)
	}
	if fmt.Sprintf("%s", err) != "Sorry, we are having trouble requesting your build from Learn. Please try again" {
		t.Errorf("error should specify that something is wrong requesting the build from learn")
	}

	// verify that requests were made properly
	if len(mockClient.Requests) != 2 {
		t.Errorf("fetching the block should make two requests")
		return
	}

	req := mockClient.Requests[0]
	if req.Method != "GET" {
		t.Errorf("Request made to Learn should be a GET, was %s", req.Method)
	}

	req = mockClient.Requests[1]
	if req.Method != "GET" {
		t.Errorf("Request made to Learn should be a GET, was %s", req.Method)
	}

	urlTarget := "https://example.com/api/v1/releases/1/release_polling?context=DIRECTORY"
	if req.URL.String() != urlTarget {
		t.Errorf("Request made to Learn should be to url '%s' but was '%s'\n", urlTarget, req.URL.String())
	}
	if req.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type header should be 'application/json', was '%s'\n", req.Header.Get("Content-Type"))
	}
	if req.Header.Get("Authorization") != "Bearer apiToken" {
		t.Errorf("Authorization header should be 'Basic apiToken', was '%s'\n", req.Header.Get("Authorization"))
	}
}

func Test_BuildReleaseFromS3_Directory(t *testing.T) {
	viper.Set("api_token", "apiToken")
	mockClient := api.MockResponse(validPreviewResponse)
	API, _ := NewAPI("https://example.com", mockClient)

	previewResponse, err := API.BuildReleaseFromS3("buket", true)
	if err != nil {
		t.Errorf("error not nil: %s\n", err)
	}
	if previewResponse.ReleaseID != 1 {
		t.Errorf("Failed to properly json parse the preview response body")
	}

	// verify that requests were made properly
	if len(mockClient.Requests) != 2 {
		t.Errorf("fetching the block should make two requests")
		return
	}

	req := mockClient.Requests[0]
	if req.Method != "GET" {
		t.Errorf("Request made to Learn should be a GET, was %s", req.Method)
	}

	req = mockClient.Requests[1]
	if req.Method != "POST" {
		t.Errorf("Request made to Learn should be a POST, was %s", req.Method)
	}

	urlTarget := "https://example.com/api/v1/releases"
	if req.URL.String() != urlTarget {
		t.Errorf("Request made to Learn should be to url '%s' but was '%s'\n", urlTarget, req.URL.String())
	}

	if req.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type header should be 'application/json', was '%s'\n", req.Header.Get("Content-Type"))
	}
	if req.Header.Get("Authorization") != "Bearer apiToken" {
		t.Errorf("Authorization header should be 'Basic apiToken', was '%s'\n", req.Header.Get("Authorization"))
	}
}

func Test_BuildReleaseFromS3_notDirectory(t *testing.T) {
	viper.Set("api_token", "apiToken")
	mockClient := api.MockResponses(credentialsResponse, validPreviewResponse)
	API, _ := NewAPI("https://example.com", mockClient)

	previewResponse, err := API.BuildReleaseFromS3("buket", false)
	if err != nil {
		t.Errorf("error not nil: %s\n", err)
	}
	if previewResponse.ReleaseID != 1 {
		t.Errorf("Failed to properly json parse the preview response body")
	}

	// verify that requests were made properly
	if len(mockClient.Requests) != 2 {
		t.Errorf("fetching the block should make two requests for credentials and one for fetching block")
		return
	}

	req := mockClient.Requests[0]
	if req.Method != "GET" {
		t.Errorf("Request made to Learn should be a GET, was %s", req.Method)
	}

	req = mockClient.Requests[1]
	if req.Method != "POST" {
		t.Errorf("Request made to Learn should be a POST, was %s", req.Method)
	}

	urlTarget := "https://example.com/api/v1/content_files"
	if req.URL.String() != urlTarget {
		t.Errorf("Request made to Learn should be to url '%s' but was '%s'\n", urlTarget, req.URL.String())
	}
	if req.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type header should be 'application/json', was '%s'\n", req.Header.Get("Content-Type"))
	}
	if req.Header.Get("Authorization") != "Bearer apiToken" {
		t.Errorf("Authorization header should be 'Basic apiToken', was '%s'\n", req.Header.Get("Authorization"))
	}
}

func Test_RetrieveCredentials(t *testing.T) {
	viper.Set("api_token", "apiToken")
	mockClient := api.MockResponse(credentialsResponse)
	API, _ := NewAPI("https://example.com", mockClient)

	if API.Credentials.S3Credentials.AccessKeyID != "access_keyin" {
		t.Errorf("Error unmarshaling S3 Credentials, access_key_id ")
	}
	if API.Credentials.SecretAccessKey != "secret_keyin" {
		t.Errorf("Error unmarshaling S3 Credentials, bad secret_access_key")
	}
	if API.Credentials.KeyPrefix != "keykey's delivery service" {
		t.Errorf("Error unmarshaling S3 Credentials, bad key_prefix")
	}
	if API.Credentials.BucketName != "buqet" {
		t.Errorf("Error unmarshaling S3 Credentials, bad bucket_name")
	}
	if API.Credentials.LatestCLIVersion != "v0.9" {
		t.Errorf("Error unmarshaling LatestCliVersion, bad latest_cli_version")
	}

	// verify that requests were made properly
	if len(mockClient.Requests) != 1 {
		t.Errorf("Creating a new api should make two requests to credentials")
		return
	}
	req := mockClient.Requests[0]
	if req.Method != "GET" {
		t.Errorf("Request made to Learn should be a GET, was %s", req.Method)
	}
	urlTarget := "https://example.com/api/v1/users/learn_cli_credentials"
	if req.URL.String() != urlTarget {
		t.Errorf("Request made to Learn should be to url '%s' but was '%s'\n", urlTarget, req.URL.String())
	}
	if req.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type header should be 'application/json', was '%s'\n", req.Header.Get("Content-Type"))
	}
	if req.Header.Get("Authorization") != "Bearer apiToken" {
		t.Errorf("Authorization header should be 'Basic apiToken', was '%s'\n", req.Header.Get("Authorization"))
	}
}
