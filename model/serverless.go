package model

type Serverless struct {
	Functions map[string]Functions `json:"functions"`
	Provider  Provider             `json:"provider"`
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

type Provider struct {
	Environment map[string]string `json:"environment"`
	Name        string            `json:"name"`
	Runtime     string            `json:"runtime"`
}
