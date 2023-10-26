package main

import (
	"fmt"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/semaphore"
)

var (
	containerSemaphore = semaphore.NewWeighted(containerNum) // 定义容器数量的信号量
)

func coderunnerProcesser(c *gin.Context) {
	var req struct {
		Code string `json:"code"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无法解析 JSON 数据"})
		return
	}

	// 尝试获取信号量，如果容器数量超出限制，会阻塞直到有空闲容器
	if err := containerSemaphore.Acquire(c.Request.Context(), 1); err != nil {
		c.JSON(429, gin.H{"error": "容器数量超出限制"})
		return
	}

	// 使用通道来并发处理请求
	responseChan := make(chan Response)
	go func() {
		defer containerSemaphore.Release(1) // 在 goroutine 结束时释放信号量

		response, found := findCache(req.Code)
		if found {
			responseChan <- response
			return
		}

		// 获取唯一的 goroutine ID 作为索引
		thisContainerIndex := atomic.AddUint64(&containerIndex, 1)
		thisContainerIndex %= containerNum
		// containerName 是 default 和 index 粘起来
		containerName := defaultContainerName + fmt.Sprint(thisContainerIndex)
		//fmt.Println("containerName:", containerName)

		var err error
		response, err = coderunner(req.Code, containerName)
		insertCache(req.Code, response)

		if err != nil {
			response.Status = 500
		} else {
			response.Status = 200
		}
		responseChan <- response
	}()

	response := <-responseChan
	c.JSON(200, response)
}
