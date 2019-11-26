package learn

import "github.com/Galvanize-IT/glearn-cli/apis"

var Api *ApiClient

// API makes network API calls to Learn
type ApiClient struct {
	client  apis.Client
	token   string
	baseUrl string
}

func NewAPI(token, baseUrl string, client apisClient) *API {
	return &API{
		client:  client,
		token:   token,
		baseUrl: baseUrl,
	}
}
