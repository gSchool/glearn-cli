package github

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gSchool/glearn-cli/api"
)

var API *APIClient

const tagUrl = "https://api.github.com/repos/gSchool/glearn-cli/tags"

type tag struct {
	Name       string
	Commit     commit
	ZipballUrl string
	TarballUrl string
	Node_id    string
}

type commit struct {
	Sha string
	Url string
}

type APIClient struct {
	Client api.Client
}

func NewAPI(client api.Client) *APIClient {
	return &APIClient{client: client}
}

func (api *APIClient) GetLatestVersion() (string, error) {
	req, err := http.NewRequest(
		"GET",
		tagUrl,
		nil,
	)
	if err != nil {
		return "", err
	}

	res, err := api.Client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	var tags []tag
	err = json.NewDecoder(res.Body).Decode(&tags)
	if err != nil {
		return "", err
	}

	fmt.Println(tags)
	if len(tags) > 0 {
		return tags[0].Name, nil
	}
	return "", fmt.Errorf("No tags found at %s\n", tagUrl)
}
