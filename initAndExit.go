package main

import (
	"os"
)

func createContainerAndFiles(thisContainerName string) {

	runCommand("run", "-d", "--name", thisContainerName, "registry.cn-hangzhou.aliyuncs.com/pivot_studio/pivot_lang:latest", "tail", "-f", "/dev/null")
	os.Mkdir("./"+thisContainerName, 0777)
	runCommand("mkdir", thisContainerName+"/"+"code")
	runCommand("touch", thisContainerName+"/"+"code/Kagari.toml", thisContainerName+"/"+"code/main.pi")

	configContent := `entry = "main.pi"
project = "main"`
	runCommand("sh", "-c", "echo '"+configContent+"' > "+thisContainerName+"/"+"code/Kagari.toml")
}
