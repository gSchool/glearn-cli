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
	Name       string `json:"name"`
	Commit     commit `json:"commit"`
	ZipballUrl string `json:"zipball_url"`
	TarballUrl string `json:"tarball_url"`
	NodeId     string `json:"node_id"`
}

type commit struct {
	Sha string `json:"sha"`
	Url string `json:"url"`
}

type APIClient struct {
	Client api.Client
}

func NewAPI(client api.Client) *APIClient {
	return &APIClient{Client: client}
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

	for _, t := range tags {
		fmt.Println(t)
	}
	if len(tags) > 0 {
		return tags[0].Name, nil
	}
	return "", fmt.Errorf("No tags found at %s\n", tagUrl)
}
