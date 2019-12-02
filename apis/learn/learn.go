package learn

import "github.com/Galvanize-IT/glearn-cli/apis"

// API is the exported APIClient, it is set during Init
var API *APIClient

// APIClient makes network API calls to Learn
type APIClient struct {
	client  apis.Client
	token   string
	baseURL string
}

// NewAPI is a constructor for the ApiClient
func NewAPI(token, baseURL string, client apis.Client) *APIClient {
	return &APIClient{
		client:  client,
		token:   token,
		baseURL: baseURL,
	}
}
