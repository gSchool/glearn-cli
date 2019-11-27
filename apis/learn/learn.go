package learn

import "github.com/Galvanize-IT/glearn-cli/apis"

// Api is the exported ApiClient, it is set during Init
var Api *ApiClient

// API makes network API calls to Learn
type ApiClient struct {
	client  apis.Client
	token   string
	baseUrl string
}

// NewAPI is a constructor for the ApiClient
func NewAPI(token, baseUrl string, client apis.Client) *ApiClient {
	return &ApiClient{
		client:  client,
		token:   token,
		baseUrl: baseUrl,
	}
}
