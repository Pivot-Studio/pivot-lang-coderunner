package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

func runCommand(command string, args ...string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(command, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func coderunner(code string, idx string) (Response, error) {

	runCommand("rm", idx+"/"+"code/main.pi")
	runCommand("touch", idx+"/"+"code/main.pi")
	runCommand("rm", "./out")

	Response := Response{}
	runCommand("sh", "-c", "echo '"+code+"' > "+idx+"/"+"code/main.pi")

	var compileoutBytes bytes.Buffer
	cmd := exec.Command("plc", idx+"/"+"code/main.pi", "-o", idx+"-out")
	cmd.Stderr = &compileoutBytes
	err := cmd.Run()
	Response.CompileOutput = compileoutBytes.String()
	if err != nil {
		fmt.Println(err)
		return Response, err
	}

	var runBytes bytes.Buffer
	cmd = exec.Command("./" + idx + "-out")
	cmd.Stdout = &runBytes
	_ = cmd.Run()
	Response.RunOutput = runBytes.String()

	return Response, nil
}
