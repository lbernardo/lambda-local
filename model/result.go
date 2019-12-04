package model

type Result struct {
	StatusCode        int               `json:"statusCode"`
	Headers           map[string]string `json:"headers,omitempty"`
	MultiValueHeaders string            `json:"multiValueHeaders,omitempty"`
	Body              string            `json:"body"`
}
