package github

import (
	"testing"

	"github.com/gSchool/glearn-cli/api"
)

const tagResponse = `[{"name":"v0.1.2"},{"name":"v0.1.1"}]`

func Test_GetLatestVersion(t *testing.T) {
	mockClient := api.MockResponse(tagResponse)
	githubClient := NewAPI(mockClient)
	version, err := githubClient.GetLatestVersion()
	if err != nil {
		t.Errorf("GetLatestVersion error: %s\n", err)
	}
	if version != "v0.1.2" {
		t.Errorf("Improper version response from mock, got '%s', expected '%s'\n", version, "v0.1.2")
	}

	req := mockClient.Requests[0]
	if req.Method != "GET" {
		t.Errorf("Request made to github should be a GET, was %s", req.Method)
	}

	urlTarget := "https://api.github.com/repos/gSchool/glearn-cli/tags"
	if req.URL.String() != urlTarget {
		t.Errorf("Request made to github should be to url '%s' but was '%s'\n", urlTarget, req.URL.String())
	}
}
