package learn

import (
	"testing"

	"github.com/gSchool/glearn-cli/api"
)

const validBlockResponse = `{"blocks":[{"id":1,"repo_name":"blocks-test","sync_errors":["somethin is wrong"],"title":"Blocks Test","cohorts_using":[7,9]}]}`

func Test_PollForBuildResponse(t *testing.T) {
	mockClient := api.MockResponse(validBlockResponse)
	API = NewAPI("apiToken", "https://example.com", mockClient)

	block, err := API.GetBlockByRepoName("blocks-test")
	if err != nil {
		t.Errorf("error not nil: %s\n", err)
	}
	testValidBlockSerialization(block, t)

	// verify that requests were made properly
	if len(mockClient.Requests) != 1 {
		t.Errorf("fetching the block should make one request")
		return
	}
	req := mockClient.Requests[0]
	if req.Method != "GET" {
		t.Errorf("Request made to Learn should be a GET, was %s", req.Method)
	}
	if req.URL.String() != "https://example.com/api/v1/blocks?repo_name=blocks-test" {
		t.Errorf("Request made to Learn should be to url '%s' but was '%s'\n", "https://example.com/api/v1/blocks?repo_name=blocks-test", req.URL.String())
	}
	if req.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type header should be 'application/json', was '%s'\n", req.Header.Get("Content-Type"))
	}
	if req.Header.Get("Authorization") != "Bearer apiToken" {
		t.Errorf("Authorization header should be 'Basic apiToken', was '%s'\n", req.Header.Get("Authorization"))
	}
}
