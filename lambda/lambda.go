package lambda

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/lbernardo/lambda-local/model"
)

func PullImageDocker(runtime string) {
	fmt.Printf("Prepare image docker")
	imageName := "lambci/lambda:" + runtime
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	_, err = cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println(" [OK]")
}

func ExecuteDockerLambda(volume string, handler string, runtime string) (model.Result, string) {
	var result model.Result
	var output bytes.Buffer

	imageName := "lambci/lambda:" + runtime

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Cmd:   []string{handler},
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

	io.Copy(&output, reader)

	str := output.String()
	str = strings.ReplaceAll(str, "\x01", "")
	str = strings.ReplaceAll(str, "\x00", "")
	str = strings.ReplaceAll(str, "J{", "{")
	str = strings.ReplaceAll(str, "\n", "")

	json.Unmarshal([]byte(str), &result)

	return result, str
}
