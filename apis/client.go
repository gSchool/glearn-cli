package apis

import (
	"bytes"
	"net/http"
)

// Client is used with http.Client and MockClient to allow mocking of services
type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

// MockClient is responsible for stubbing requests in tests
// if the Response field is set, the mock client will respond to each request with the Response
// If Responses is set, the mock client will match subsequent requests to subsequent responses,
// moving along the array of Responses once for each request
type MockClient struct {
	Response   []byte
	Responses  [][]byte
	Requests   []*http.Request
	StatusCode int
}

// Do saves the HTTP Request, returning 200 and no error
func (mock *MockClient) Do(req *http.Request) (*http.Response, error) {
	mock.Requests = append(mock.Requests, req)
	statusCode := 200
	if mock.StatusCode != 0 && mock.StatusCode != 200 {
		statusCode = mock.StatusCode
	}
	if len(mock.Response) > 0 {
		return &http.Response{
			StatusCode: statusCode,
			Body:       MockBody(mock.Response),
		}, nil
	}

	if len(mock.Responses) > 0 {
		return &http.Response{
			StatusCode: statusCode,
			Body:       MockBody(mock.Responses[len(mock.Requests)-1]),
		}, nil
	}
	return &http.Response{StatusCode: 200}, nil
}

// MockResponses creates a mock client with a series of canned responses
func MockResponses(bodies ...string) *MockClient {
	mock := &MockClient{}
	for _, body := range bodies {
		mock.Responses = append(mock.Responses, []byte(body))
	}
	return mock
}

// MockResponse creates a mock client that will return the given response on the body
func MockResponse(body string) *MockClient {
	mock := &MockClient{}
	mock.Response = []byte(body)
	return mock
}

// MockBody creates a mock body that can be added to a response
func MockBody(b []byte) *ClosingBuffer {
	return &ClosingBuffer{bytes.NewBuffer(b)}
}

// ClosingBuffer embeds a bytes.Buffer, giving it a Read method required for a Response Body ReadCloser body
type ClosingBuffer struct {
	*bytes.Buffer
}

// Close allows a ClosingBuffer to Close, implementing a net/http Response Body ReadCloser interface
func (cb ClosingBuffer) Close() (err error) {
	return
}
