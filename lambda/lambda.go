package lambda

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/lbernardo/lambda-local/model"
)

func ExecuteDockerLambda(volume string, handler string, runtime string) model.Result {
	var result model.Result
	var out bytes.Buffer
	var out2 bytes.Buffer

	imageName := "lambci/lambda:" + runtime
	var content model.LambdaContent
	var responseCreate model.CreateResponse

	content.Image = imageName
	content.Cmd = []string{handler}
	content.HostConfig.Binds = []string{volume + ":/var/task"}

	body, _ := content.Marshal()

	cmd := exec.Command("curl", "--unix-socket", "/var/run/docker.sock", "-H", "Content-Type: application/json", "-d", string(body), "-X", "POST", "http:/v1.24/containers/create")
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	cmd.Run()

	json.Unmarshal(out.Bytes(), &responseCreate)

	cmd = exec.Command("curl", "--unix-socket", "/var/run/docker.sock", "-X", "POST", "http:/v1.24/containers/"+responseCreate.ID+"/start")
	cmd.Run()

	cmd = exec.Command("curl", "--unix-socket", "/var/run/docker.sock", "-s", "-o", "-", "http:/v1.24/containers/"+responseCreate.ID+"/logs?stdout=1")
	cmd.Stdout = &out2
	cmd.Stderr = os.Stderr
	cmd.Run()

	json.Unmarshal(out2.Bytes(), &result)
	fmt.Println(out2.String())
	fmt.Println(result)

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
