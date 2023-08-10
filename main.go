package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	containerName = "coderunner"
)

func main() {
	r := gin.Default()
	r.POST("/coderunner", coderunnerProcesser)
	r.Run(":8080")
}

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

func coderunnerProcesser(c *gin.Context) {
	createContainerAndFiles()
	var req struct {
		Code string `json:"code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无法解析 JSON 数据"})
		return
	}

	result, err := coderunner(req.Code)
	if err != nil {
		c.JSON(200, gin.H{"result": err.Error()})
		return
	}
	c.JSON(200, gin.H{"result": result})
}

func coderunner(code string) (string, error) {
	runDockerCommand("exec", containerName, "sh", "-c", "echo '"+code+"' >> code/main.pi")
	runDockerCommand("exec", containerName, "plc", "code/main.pi")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("docker", "exec", containerName, "./out")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	runDockerCommand("exec", containerName, "rm", "code/main.pi")
	runDockerCommand("exec", containerName, "touch", "code/main.pi")
	runDockerCommand("exec", containerName, "rm", "./out")

	if err != nil {
		return "", fmt.Errorf("运行时出错：%s\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}
