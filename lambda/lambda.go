package lambda

import (
	"os/exec"

	"github.com/lbernardo/lambda-local/model"
)

func ExecuteDockerLambda(volume string, handler string, runtime string) model.Result {
	var result model.Result

	imageName := "lambci/lambda:" + runtime
	var content model.LambdaContent

	content.Image = imageName
	content.Cmd = []string{handler}
	content.HostConfig.Binds = []string{volume + ":/var/task"}

	body, _ := content.Marshal()

	cmd := exec.Command("curl", "--unix-socket", "/var/run/docker.sock", "-H", "Content-Type: application/json", "-d", string(body), "-X", "POST", "http:/v1.24/containers/create")
	cmd.Run()

	cmd := exec.Command()

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
