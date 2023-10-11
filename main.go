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
	semaphore      = make(chan struct{}, containerNum) // 信号量，控制同时处理的请求数量
	containerIndex = 0
	containerName  = ""
	imageName      = "registry.cn-hangzhou.aliyuncs.com/pivot_studio/pivot_lang"
	imageTag       = "latest"
	updateInterval = 7 * 24 * time.Hour
)

func init() {
	createCache()
	for i := 0; i < containerNum; i++ {
		semaphore <- struct{}{}
	}
}

func main() {
	go updateDockerImage(imageName, imageTag, updateInterval) // 启动 Docker 镜像更新任务

	r := gin.Default()
	r.POST("/coderunner", coderunnerProcesser)
	r.Run(":8080")

	log.Fatal(autotls.Run(r, "code.lang.pivotstudio.cn"))
}
