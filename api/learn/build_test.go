package learn

import (
	"testing"

	"github.com/gSchool/glearn-cli/api"
	"github.com/gSchool/glearn-cli/api/learn"
)

const validBlockRespons = `{"blocks":[{"id":1,"repo_name":"blocks-test","sync_errors":["somethin is wrong"],"title":"Blocks Test","cohorts_using":[7,9]}]}`

func Test_GetBlockByRepoName(t *testing.T) {
	mockClient := api.MockResponse(validBlock)
	learn.API = learn.NewAPI("apiToken", "http://example.com", &mockClient)

	block, err := learn.API.CreateBlockByRepoName("repo-name")
	if block.ID != 1 {
		t.Errorf("block response should have id of 1, but got %d\n", block.id)
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
