package main

import (
	"log"

	"time"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
)

const (
	defaultContainerName = "coderunner"
	cacheTTS             = 10
	containerNum         = 10
)

type Response struct {
	Status        int    `json:"status"`
	CompileOutput string `json:"compileOutput"`
	RunOutput     string `json:"runOutput"`
}

var (
	containerIndex = uint64(0)
	imageName      = "registry.cn-hangzhou.aliyuncs.com/pivot_studio/pivot_lang"
	imageTag       = "latest"
	updateInterval = 7 * 24 * time.Hour
)

func init() {
	createCache()
}

func main() {
	go updateDockerImage(imageName, imageTag, updateInterval) // 启动 Docker 镜像更新任务

	r := gin.Default()
	r.POST("/coderunner", coderunnerProcesser)
	r.Run(":8080")

	log.Fatal(autotls.Run(r, "code.lang.pivotstudio.cn"))

	defer deleteAllContainers()
}
