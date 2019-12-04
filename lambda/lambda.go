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
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/lbernardo/lambda-local/model"
)

type ContentRequest struct {
	Body string `json:"body"`
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

func ExecuteDockerLambda(volume string, handler string, runtime string, body io.ReadCloser) (model.Result, string) {
	var result model.Result
	var output bytes.Buffer
	var contentRequest ContentRequest
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)
	bodyStr := buf.String()

	imageName := "lambci/lambda:" + runtime

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	bodyStr = strings.ReplaceAll(bodyStr, "\t", "")
	bodyStr = strings.ReplaceAll(bodyStr, "\n", "")
	contentRequest.Body = bodyStr

	jsonRequest, _ := json.Marshal(contentRequest)

	var executeCommand []string
	executeCommand = append(executeCommand, handler)
	executeCommand = append(executeCommand, string(jsonRequest))

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Cmd:   executeCommand,
	}, &container.HostConfig{
		Binds: []string{volume + ":/var/task"},
	}, nil, "")
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
