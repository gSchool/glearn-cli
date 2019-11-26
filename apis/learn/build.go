package learn

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Block struct {
	id            int
	repo_name     string
	sync_errors   []string
	title         string
	cohorts_using []int
}

type blockResponse struct {
	blocks []Block `json:"blocks"`
}

func (b Block) Exists() bool {
	return b.id != 0
}

func (api *ApiClient) GetBlockByRepoName(repoName string) (block, error) {
	u, err := url.Parse(fmt.Sprintf("%s/api/v1/blocks", api.baseUrl))
	if err != nil {
		return block{}, errors.New("unable to parse Learn remote")
	}
	v := url.Values{}
	v.Set("repo_name", repoName)
	u.RawQuery = v.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
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

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error: response status: %d", res.StatusCode)
	}

	var blockResp blockResponse
	json.NewDecoder(res.Body).Decode(responseBody)

	if len(responseBody.blocks) == 1 {
		return responseBody.blocks[0]
	}
	return &Block, nil
}

func (api *ApiClient) CreateBlockByRepoName(repoName string) (block, error) {
	payload := map[string]string{
		"bucket_key_name": bucketKey,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: time.Second * 10}

	req, err := http.NewRequest("POST", "https://httpbin.org/post", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api.token))

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error: response status: %d", res.StatusCode)
	}

}
