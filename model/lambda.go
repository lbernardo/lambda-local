package model

import "encoding/json"

type LambdaContent struct {
	Image      string     `json:"Image"`
	Cmd        []string   `json:"Cmd"`
	HostConfig HostConfig `json:"HostConfig"`
}

type HostConfig struct {
	Binds []string `json:"Binds"`
}

func (r *LambdaContent) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
