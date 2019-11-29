package lambda

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/lbernardo/lambda-local/model"
)

func ExecuteDockerLambda(volume string, handler string, runtime string) model.Result {
	var result model.Result

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

	p := make([]byte, 8)
	reader.Read(p)
	content, _ := ioutil.ReadAll(reader)
	fmt.Println(string(content))
	return result
}
