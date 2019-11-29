package lambda

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/lbernardo/lambda-local/model"
)

func ExecuteDockerLambda(volume string, handler string, runtime string) model.Result {
	var result model.Result
	// var out bytes.Buffer
	// var out2 bytes.Buffer

	imageName := "lambci/lambda:" + runtime

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, out)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	// cmd := exec.Command("docker", "run", "--rm", "-v", volume+":/var/task", "lambci/lambda:"+runtime, handler)
	// cmd.Stdout = &out
	// cmd.Stderr = os.Stderr

	// err := cmd.Run()

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// json.Unmarshal(out.Bytes(), &result)

	return result
}
