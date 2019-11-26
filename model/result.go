package model

type Result struct {
	StatusCode        int    `json:"statusCode"`
	Headers           string `json:"headers"`
	MultiValueHeaders string `json:"multiValueHeaders"`
	Body              string `json:"body"`
}
