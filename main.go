package main

import (
	"github.com/gin-gonic/gin"
)

const (
	containerName = "coderunner"
	cacheTTS      = 10
)

type Response struct {
	Status        int    `json:"status"`
	CompileOutput string `json:"compileOutput"`
	RunOutput     string `json:"runOutput"`
}

func main() {
	createCache()
	createContainerAndFiles()
	r := gin.Default()
	r.POST("/coderunner", coderunnerProcesser)
	r.Run(":8080")
}
