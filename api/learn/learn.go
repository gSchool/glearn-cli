package learn

import "github.com/gSchool/glearn-cli/api"

// API is the exported APIClient, it is set during Init
var API *APIClient

// APIClient makes network API calls to Learn
type APIClient struct {
	client  api.Client
	token   string
	baseURL string
}

// NewAPI is a constructor for the ApiClient
func NewAPI(token, baseURL string, client api.Client) *APIClient {
	return &APIClient{
		client:  client,
		token:   token,
		baseURL: baseURL,
	}
}
