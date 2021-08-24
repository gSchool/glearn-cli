package github

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func Test_GetLatestVersion(t *testing.T) {
	client := &http.Client{Timeout: 15 * time.Second}
	githubClient := NewAPI(client)
	version, err := githubClient.GetLatestVersion()
	if err != nil {
		t.Errorf("GetLatestVersion error: %s\n", err)
	}
	fmt.Println(version)
}
