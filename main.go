package main

import (
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
)

func init() {
	createCache()
	for i := 0; i < containerNum; i++ {
		semaphore <- struct{}{}
	}
}

func main() {
	r := gin.Default()
	r.POST("/coderunner", coderunnerProcesser)
	r.Run(":8080")
}
