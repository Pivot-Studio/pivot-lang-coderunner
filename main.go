package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

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
)

func init() {
	createCache()
}

func main() {
	// go updateDockerImage(imageName, imageTag, updateInterval) // 启动 Docker 镜像更新任务

	r := gin.Default()
	r.Use(CORSMiddleware())
	r.POST("/coderunner", coderunnerProcesser)
	log.Fatal(r.Run(":8080"))

	// log.Fatal(autotls.Run(r, "code.lang.pivotstudio.cn"))

}
