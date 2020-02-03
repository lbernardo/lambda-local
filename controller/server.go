package controller

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/lbernardo/lambda-local/lambda"
	"github.com/lbernardo/lambda-local/model"
)

type Server struct {
	Host            string
	Port            string
	Yaml            string
	Volume          string
	Network         string
	EnvironmentFile string
	JSON            model.Serverless
}

func (se *Server) ContentYaml() {
	content, err := ReadYaml(se.Yaml)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(content, &se.JSON)
}

func (se *Server) ReadEnv() {
	se.JSON.Provider.Environment = make(map[string]string, 0)
	if se.EnvironmentFile != "" {
		file, err := os.Open(se.EnvironmentFile)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			textLine := scanner.Text()
			if len(textLine) <= 0 {
				continue
			}
			envArgs := strings.Split(textLine, "=")
			se.JSON.Provider.Environment[envArgs[0]] = envArgs[1]
		}
	}
}

func (se *Server) StartConfig() {
	se.ContentYaml()
	se.ReadEnv()
	lambda.PullImageDocker(se.JSON.Provider.Runtime)

	route := mux.NewRouter()

	for _, functions := range se.JSON.Functions {
		fmt.Println(functions.Events[0].HttpEvent.Path)
		route.HandleFunc("/"+functions.Events[0].HttpEvent.Path, func(w http.ResponseWriter, r *http.Request) {
			parameters := mux.Vars(r)
			fmt.Println(parameters)
			fmt.Println(functions.Events[0].HttpEvent.Path)

			result, off := lambda.ExecuteDockerLambda(se.Volume, se.Network, functions.Handler, se.JSON.Provider.Runtime, se.JSON.Provider.Environment, r.Body, parameters)
			if result.StatusCode == 0 {
				w.WriteHeader(400)
				fmt.Println(off)
				return
			}

			for key, val := range result.Headers {
				w.Header().Set(key, val)
			}
			w.WriteHeader(result.StatusCode)
			w.Write([]byte(result.Body))
			return
		}).Methods(functions.Events[0].HttpEvent.Method)
	}

	http.Handle("/", route)
}

func (se *Server) StartServer() {
	fmt.Printf("Start server API Gateway -> lambda %v:%v\n", se.Host, se.Port)
	log.Fatal(http.ListenAndServe(se.Host+":"+se.Port, nil))
}
