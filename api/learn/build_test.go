package learn

import (
	"testing"

	"github.com/gSchool/glearn-cli/api"
	"github.com/spf13/viper"
)

const validBlockResponse = `{"blocks":[{"id":1,"repo_name":"blocks-test","sync_errors":["somethin is wrong"],"title":"Blocks Test","cohorts_using":[7,9]}]}`

func Test_GetBlockByRepoName(t *testing.T) {
	viper.Set("api_token", "apiToken")
	mockClient := api.MockResponse(validBlockResponse)
	API, _ := NewAPI("https://example.com", mockClient)

	block, err := API.GetBlockByRepoName("blocks-test")
	if err != nil {
		t.Errorf("error not nil: %s\n", err)
	}
	testValidBlockSerialization(block, t)

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

func Test_CreateBlockByRepoName(t *testing.T) {
	viper.Set("api_token", "apiToken")
	mockClient := api.MockResponse(validBlockResponse)
	API, _ := NewAPI("https://example.com", mockClient)

	block, err := API.CreateBlockByRepoName("blocks-test")
	if err != nil {
		t.Errorf("error not nil: %s\n", err)
	}
	testValidBlockSerialization(block, t)

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
	if req.URL.String() != "https://example.com/api/v1/blocks" {
		t.Errorf("Request made to Learn should be to url '%s' but was '%s'\n", "https://example.com/api/v1/blocks", req.URL.String())
	}
	if req.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type header should be 'application/json', was '%s'\n", req.Header.Get("Content-Type"))
	}
	if req.Header.Get("Authorization") != "Bearer apiToken" {
		t.Errorf("Authorization header should be 'Basic apiToken', was '%s'\n", req.Header.Get("Authorization"))
	}
}

const validMasterReleaseResponse = `{"release_id":9}`

func Test_CreateMasterRelease(t *testing.T) {
	viper.Set("api_token", "apiToken")
	mockClient := api.MockResponse(validMasterReleaseResponse)
	API, _ := NewAPI("https://example.com", mockClient)

	id, err := API.CreateMasterRelease(1)
	if err != nil {
		t.Errorf("error not nil: %s\n", err)
	}
	if id != 9 {
		t.Errorf("Response release id was %d but expected 9", id)
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
	if req.URL.String() != "https://example.com/api/v1/blocks/1/releases" {
		t.Errorf("Request made to Learn should be to url '%s' but was '%s'\n", "https://example.com/api/v1/blocks/1/releases", req.URL.String())
	}
	if req.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Content-Type header should be 'application/json', was '%s'\n", req.Header.Get("Content-Type"))
	}
	if req.Header.Get("Authorization") != "Bearer apiToken" {
		t.Errorf("Authorization header should be 'Basic apiToken', was '%s'\n", req.Header.Get("Authorization"))
	}
}

func testValidBlockSerialization(block Block, t *testing.T) {
	if block.ID != 1 {
		t.Errorf("block response should have id of 1, but got %d\n", block.ID)
	}
	if block.RepoName != "blocks-test" {
		t.Errorf("block response should have repo_name of 'blocks-test', but got %s\n", block.RepoName)
	}
	if len(block.SyncErrors) != 1 && block.SyncErrors[0] != "somethin is wrong" {
		t.Errorf("block response should have sync_errors of ['somethin is wrong'], but got %+v\n", block.SyncErrors)
	}
	if block.Title != "Blocks Test" {
		t.Errorf("block response should have title of 'Blocks Test', but got %s\n", block.Title)
	}
	if len(block.CohortsUsing) != 2 && block.CohortsUsing[0] != 7 && block.CohortsUsing[1] != 9 {
		t.Errorf("block response should have cohorts_using of [7,9], but got %+v\n", block.CohortsUsing)
	}
}
