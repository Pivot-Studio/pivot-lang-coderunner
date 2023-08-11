package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
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

func createContainerAndFiles() {
	stdout, _, err := runDockerCommand("ps", "-a", "--filter", "name="+containerName)
	if err != nil || strings.Contains(stdout, containerName) {
		runDockerCommand("start", containerName)
		return
	}

	runDockerCommand("run", "-d", "--name", containerName, "registry.cn-hangzhou.aliyuncs.com/pivot_studio/pivot_lang:latest", "tail", "-f", "/dev/null")
	runDockerCommand("exec", containerName, "mkdir", "code")
	runDockerCommand("exec", containerName, "touch", "code/Kagari.toml", "code/main.pi")

	configContent := `entry = "main.pi"
project = "main"`
	runDockerCommand("exec", containerName, "sh", "-c", "echo '"+configContent+"' >> code/Kagari.toml")
}

func coderunner(code string) (Response, error) {

	runDockerCommand("exec", containerName, "rm", "code/main.pi")
	runDockerCommand("exec", containerName, "touch", "code/main.pi")
	runDockerCommand("exec", containerName, "rm", "./out")

	Response := Response{}
	runDockerCommand("exec", containerName, "sh", "-c", "echo '"+code+"' >> code/main.pi")

	var compileoutBytes bytes.Buffer
	cmd := exec.Command("docker", "exec", containerName, "plc", "code/main.pi")
	cmd.Stdout = &compileoutBytes
	err := cmd.Run()
	Response.CompileOutput = compileoutBytes.String()
	if err != nil {
		fmt.Println(err)
		return Response, err
	}

	var runBytes bytes.Buffer
	cmd = exec.Command("docker", "exec", containerName, "./out")
	cmd.Stdout = &runBytes
	_ = cmd.Run()
	Response.RunOutput = runBytes.String()
	fmt.Println(Response.RunOutput)

	runDockerCommand("exec", containerName, "rm", "code/main.pi")
	runDockerCommand("exec", containerName, "touch", "code/main.pi")
	runDockerCommand("exec", containerName, "rm", "./out")

	return Response, nil
}
