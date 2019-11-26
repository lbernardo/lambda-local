package lambda

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"os/exec"

	"github.com/lbernardo/lambda-local/model"
)

func ExecuteDockerLambda(volume string, handler string, runtime string) model.Result {
	var out bytes.Buffer
	var result model.Result

	cmd := exec.Command("docker", "run", "--rm", "-v", volume+":/var/task", "lambci/lambda:"+runtime, handler)
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(out.Bytes(), &result)

	return result
}
