package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/lbernardo/lambda-local/lambda"
	"github.com/lbernardo/lambda-local/model"
)

type Server struct {
	Host   string
	Port   string
	Yaml   string
	Volume string
	JSON   model.Serverless
}

func checkPath(event model.HttpEvent, reqPath string, method string) bool {
	emethod := strings.ToUpper(event.Method)

	event.Path = strings.ReplaceAll("/"+event.Path, "//", "/")

	if emethod != method {
		return false
	}

	if event.Path == reqPath {
		return true
	}

	match, _ := regexp.MatchString("{[a-z0-9A-Z-]+}", event.Path)
	if match == true {
		reg, _ := regexp.Compile("{[a-z0-9A-Z-]+}")
		ep := reg.ReplaceAllString(event.Path, "[a-z0-9A-Z-]+")
		match, _ = regexp.MatchString(ep, reqPath)
		if match == true {
			return true
		}
	}

	return false
}

func (se *Server) StartConfig() {
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content, err := ReadYaml(se.Yaml)
		if err != nil {
			panic(err)
		}
		json.Unmarshal(content, &se.JSON)
		for _, functions := range se.JSON.Functions {
			if checkPath(functions.Events[0].HttpEvent, r.URL.RequestURI(), r.Method) {
				result := lambda.ExecuteDockerLambda(se.Volume, functions.Handler, se.JSON.Provider["runtime"])
				w.WriteHeader(result.StatusCode)
				w.Write([]byte(result.Body))
				return
			}
		}
		w.WriteHeader(404)
		w.Write([]byte("404 Not found"))
	}))
}

func (se *Server) StartServer() {
	fmt.Printf("Start server lambda %v:%v\n", se.Host, se.Port)
	log.Fatal(http.ListenAndServe(se.Host+":"+se.Port, nil))
}
