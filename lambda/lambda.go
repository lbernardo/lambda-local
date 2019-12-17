package lambda

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/lbernardo/lambda-local/model"
)

type ContentRequest struct {
	Body           string            `json:"body"`
	PathParameters map[string]string `json:"pathparameters"`
}

func PullImageDocker(runtime string) {
	fmt.Println("Prepare image docker")
	imageName := "lambci/lambda:" + runtime
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	reader, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)
}

func ReplaceEnvironment(env string) string {
	return strings.ReplaceAll(env, "${opt:stage, self:provider.stage}", "dev")
}

func ExecuteDockerLambda(volume string, net string, handler string, runtime string, environment map[string]string, body io.ReadCloser, parameters map[string]string) (model.Result, string) {
	var result model.Result
	var output bytes.Buffer
	var contentRequest ContentRequest
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)
	bodyStr := buf.String()
	var strEnv []string

	imageName := "lambci/lambda:" + runtime

	for n, env := range environment {
		strEnv = append(strEnv, n+"="+ReplaceEnvironment(env))
	}

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	bodyStr = strings.ReplaceAll(bodyStr, "\t", "")
	bodyStr = strings.ReplaceAll(bodyStr, "\n", "")
	contentRequest.Body = bodyStr
	contentRequest.PathParameters = parameters

	jsonRequest, _ := json.Marshal(contentRequest)

	var executeCommand []string
	executeCommand = append(executeCommand, handler)
	executeCommand = append(executeCommand, string(jsonRequest))

	// Network config
	networkingConfig := &network.NetworkingConfig{}

	if net != "" {
		networkingConfig.EndpointsConfig = map[string]*network.EndpointSettings{
			net: &network.EndpointSettings{},
		}
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Cmd:   executeCommand,
		Env:   strEnv,
	}, &container.HostConfig{
		Binds: []string{volume + ":/var/task"},
	}, networkingConfig, "")
	if err != nil {
		panic(err)
	}

	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})

	if err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	reader, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(&output, os.Stderr, reader)

	str := output.String()
	err = json.Unmarshal([]byte(str), &result)
	if err != nil {
		fmt.Println(err)
	}

	return result, str
}
