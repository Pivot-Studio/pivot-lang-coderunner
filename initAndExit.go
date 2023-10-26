package main

import (
	"fmt"
	"strings"
)

func createContainerAndFiles(thisContainerName string) {
	stdout, _, err := runDockerCommand("ps", "-a", "--filter", "name="+thisContainerName)
	if err != nil || strings.Contains(stdout, thisContainerName) {
		runDockerCommand("start", thisContainerName)
		//加一个资源限制
		runDockerCommand("update", "--cpus", "0.5", "--memory", "512m", thisContainerName)
		return
	}

	runDockerCommand("run", "-d", "--name", thisContainerName, "registry.cn-hangzhou.aliyuncs.com/pivot_studio/pivot_lang:latest", "tail", "-f", "/dev/null")
	runDockerCommand("exec", thisContainerName, "mkdir", "code")
	runDockerCommand("exec", thisContainerName, "touch", "code/Kagari.toml", "code/main.pi")

	configContent := `entry = "main.pi"
project = "main"`
	runDockerCommand("exec", thisContainerName, "sh", "-c", "echo '"+configContent+"' >> code/Kagari.toml")
}

func deleteAllContainers() {
	for i := 1; i < containerNum+1; i++ {
		thisContainerName := defaultContainerName + "_" + fmt.Sprint(i)
		runDockerCommand("stop", thisContainerName)
		runDockerCommand("rm", thisContainerName)
	}
}
