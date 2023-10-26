package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func runDockerCommand(args ...string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("docker", args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func coderunner(code string, thisContainerName string) (Response, error) {

	runDockerCommand("exec", thisContainerName, "rm", "code/main.pi")
	runDockerCommand("exec", thisContainerName, "touch", "code/main.pi")
	runDockerCommand("exec", thisContainerName, "rm", "./out")

	Response := Response{}
	runDockerCommand("exec", thisContainerName, "sh", "-c", "echo '"+code+"' >> code/main.pi")

	var compileoutBytes bytes.Buffer
	cmd := exec.Command("docker", "exec", thisContainerName, "plc", "code/main.pi")
	cmd.Stderr = &compileoutBytes
	err := cmd.Run()
	Response.CompileOutput = compileoutBytes.String()
	if err != nil {
		fmt.Println(err)
		//delete container
		runDockerCommand("stop", thisContainerName)
		runDockerCommand("rm", thisContainerName)
		return Response, err
	}

	var runBytes bytes.Buffer
	cmd = exec.Command("docker", "exec", thisContainerName, "./out")
	cmd.Stdout = &runBytes
	_ = cmd.Run()
	Response.RunOutput = runBytes.String()

	//delete container
	runDockerCommand("stop", thisContainerName)
	runDockerCommand("rm", thisContainerName)

	return Response, nil
}
