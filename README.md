# lambda-local
Execute lambda aws functions local


```bash
$ lambda-local 
Execute functions lambda aws

Usage:
  lambda-local [command]

Available Commands:
  help        Help about any command
  start       Start local functions lambda

Flags:
  -h, --help   help for lambda-local

Use "lambda-local [command] --help" for more information about a command.
subcommand is required

```

Start local lambda functions development


## Build
For build
```
make build
```
or
```
go build -o lambda-local
```

## Install
```
go install github.com/lbernardo/lambda-local
```

**Note: $GOPATH/bin must be set to $PATH**


## Start local functions lambda

```
$ lambda-local start --volume $PWD

Start server lambda 0.0.0.0:3000
```

```
Start local functions lambda

Usage:
  lambda-local start [flags]

Flags:
  -h, --help            help for start
      --host string     host usage [default 0.0.0.0] (default "0.0.0.0")
      --port string     port usage [default 3000] (default "3000")
      --volume string   Docker volume mount execution [required] [ [ex: --volume $PWD]
      --yaml string     File yaml serverless [default serverless.yml]
```

## Example usage

**serverless.yml**

```
service: myservice

provider:
  name: aws
  runtime: go1.x

package:
 exclude:
   - ./**
 include:
   - ./bin/**

functions:
  hello:
    handler: bin/main
    events:
      - http:
          path: hello
          method: get
```

**Makefile**
```make
build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/hello adapters/hello/*
```

**adapter/hello/hello.go**

```go
package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       "Hello",
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}

```


**terminal**
```
$ make build

> env GOOS=linux go build -ldflags="-s -w" -o bin/hello adapters/hello/* 

$ lambda-local start --volume $PWD --port 5000

> Start server lambda 0.0.0.0:5000
```

**curl**
```
$ curl http://localhost:5000/hello
> Hello
```
