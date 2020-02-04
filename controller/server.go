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

func (se *Server) executeHandler(handler string, ) {

}

func (se *Server) StartConfig() *mux.Router {
	lambda.PullImageDocker(se.JSON.Provider.Runtime)

	route := mux.NewRouter()

	for _, functions := range se.JSON.Functions {
		path := "/" + functions.Events[0].HttpEvent.Path
		method := functions.Events[0].HttpEvent.Method
		function := functions.Handler

		path = strings.ReplaceAll(path,"//","/")

		fmt.Printf("http://%v:%v%v [%v]\n\n", se.Host, se.Port, path, strings.ToUpper(method))
		fff := func(w http.ResponseWriter, r *http.Request) {
			parameters := mux.Vars(r)
			result, off := lambda.ExecuteDockerLambda(se.Volume, se.Network, function, se.JSON.Provider.Runtime, se.JSON.Provider.Environment, r.Body, parameters)
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
		}
		route.HandleFunc(path, fff).Methods(method)
	}

	return route
}

func (se *Server) StartServer() {
	se.ContentYaml()
	se.ReadEnv()

	fmt.Printf("Start server API Gateway -> lambda %v:%v\n", se.Host, se.Port)
	log.Fatal(http.ListenAndServe(se.Host+":"+se.Port, se.StartConfig()))
}
