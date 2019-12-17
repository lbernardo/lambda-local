package model

type Serverless struct {
	Functions   map[string]Functions `json:"functions"`
	Provider    map[string]string    `json:"provider"`
	Environment map[string]string    `json:"environment"`
}

type Functions struct {
	Events  []Event `json:"events"`
	Handler string  `json:"handler"`
}

type Event struct {
	HttpEvent HttpEvent `json:"http"`
}

type HttpEvent struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}
