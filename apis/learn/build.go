package learn

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// Block holds information yielded from the Learn Block API
type Block struct {
	Id           int      `json:"id"`
	RepoName     string   `json:"repo_name"`
	SyncErrors   []string `json:"sync_errors"`
	Title        string   `json:"title"`
	CohortsUsing []int    `json:"cohorts_using"`
}

type blockResponse struct {
	Blocks []Block `json:"blocks"`
}

type BlockPost struct {
	Block Block `json:"block"`
}

// ReleaseResponse holds the release id of a fetched or created release
type ReleaseResponse struct {
	ReleaseID int `json:"release_id"`
}

// Exists reports if a Block struct has a nonzero id value
func (b Block) Exists() bool {
	return b.Id != 0
}

func (api *ApiClient) GetBlockByRepoName(repoName string) (Block, error) {
	u, err := url.Parse(fmt.Sprintf("%s/api/v1/blocks", api.baseUrl))
	if err != nil {
		return Block{}, errors.New("unable to parse Learn remote")
	}
	v := url.Values{}
	v.Set("repo_name", repoName)
	u.RawQuery = v.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return Block{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api.token))

	res, err := api.client.Do(req)
	if err != nil {
		return Block{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Block{}, fmt.Errorf("Error: response status: %d", res.StatusCode)
	}

	var blockResp blockResponse
	err = json.NewDecoder(res.Body).Decode(&blockResp)
	if err != nil {
		return Block{}, err
	}

	if len(blockResp.Blocks) == 1 {
		return blockResp.Blocks[0], nil
	}
	return Block{}, nil
}

func (api *ApiClient) CreateBlockByRepoName(repoName string) (Block, error) {
	payload := BlockPost{Block: Block{RepoName: repoName}}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return Block{}, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/blocks", api.baseUrl), bytes.NewBuffer(payloadBytes))
	if err != nil {
		return Block{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api.token))

	res, err := api.client.Do(req)
	if err != nil {
		return Block{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Block{}, fmt.Errorf("Error: response status: %d", res.StatusCode)
	}

	var blockResp blockResponse
	err = json.NewDecoder(res.Body).Decode(&blockResp)
	if err != nil {
		return Block{}, err
	}

	if len(blockResp.Blocks) == 1 {
		return blockResp.Blocks[0], nil
	}
	return Block{}, nil
}

func (api *ApiClient) CreateMasterRelease(blockId int) (int, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/blocks/%d/releases", api.baseUrl, blockId), nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", api.token))

	res, err := api.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Error: response status: %d", res.StatusCode)
	}
	var r ReleaseResponse
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return 0, err
	}
	return r.ReleaseID, nil
}
