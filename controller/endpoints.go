package controller

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lbernardo/lambda-local/model"
)

type EndpointsController struct {
	Yaml string
	Host string
	Port string
	JSON model.Serverless
}

func NewEndpointsController(yaml string, host string, port string) *EndpointsController {
	return &EndpointsController{
		Yaml: yaml,
		Host: host,
		Port: port,
	}
}

func (ec *EndpointsController) ContentYaml() {
	content, err := ReadYaml(ec.Yaml)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(content, &ec.JSON)
}

func (ec *EndpointsController) ListEndpoints() {
	ec.ContentYaml()
	fmt.Println("Endpoints")
	for _, function := range ec.JSON.Functions {
		for _, event := range function.Events {
			fmt.Printf("- http://%v:%v/%v [\033[01;32m%v\033[0m]\n\n", ec.Host, ec.Port, event.HttpEvent.Path, strings.ToUpper(event.HttpEvent.Method))
		}
	}
}
