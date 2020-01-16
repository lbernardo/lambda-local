package controller

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

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

var strReg = "{[a-z0-9A-Z-]+}"

const varsKey int = iota

func checkPath(event model.HttpEvent, reqPath string, method string) (bool, map[string]string) {
	parameters := map[string]string{}
	emethod := strings.ToUpper(event.Method)

	event.Path = strings.ReplaceAll("/"+event.Path, "//", "/")

	if emethod != method {
		return false, parameters
	}

	if event.Path == reqPath {
		return true, parameters
	}

	match, _ := regexp.MatchString(strReg, event.Path)
	if match == true {
		reg, _ := regexp.Compile(strReg)

		ep := reg.ReplaceAllString(event.Path, "[a-z0-9A-Z-]+")
		p := strings.ReplaceAll(ep, "[a-z0-9A-Z-]+", "")

		match, _ = regexp.MatchString(ep, reqPath)
		if match == true {
			reg2, _ := regexp.Compile("[a-z0-9A-Z-]+")
			params := reg.FindStringSubmatch(event.Path)
			m := strings.ReplaceAll(reqPath, p, "")
			values := reg2.FindStringSubmatch(m)
			var value string

			for i, param := range params {
				value = values[i]
				param = strings.ReplaceAll(param, "{", "")
				param = strings.ReplaceAll(param, "}", "")
				parameters[param] = value
			}

			return true, parameters
		}
	}

	return false, parameters
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
		file, err := os.Open(".local.env")
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
	lambda.PullImageDocker(se.JSON.Provider.Runtime)

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		se.ContentYaml()
		se.ReadEnv()
		for _, functions := range se.JSON.Functions {
			check, parameters := checkPath(functions.Events[0].HttpEvent, r.URL.RequestURI(), r.Method)
			if check {
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
