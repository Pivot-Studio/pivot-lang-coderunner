package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func coderunnerProcesser(c *gin.Context) {
	var req struct {
		Code string `json:"code"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无法解析 JSON 数据"})
		return
	}

	// 使用通道来并发处理请求
	responseChan := make(chan Response, containerNum)
	go func() {
		response, found := findCache(req.Code)
		if found {
			responseChan <- response
			return
		}
		containerIndex = containerIndex % containerNum
		containerIndex++
		//containername是default和index粘起来
		containerName = defaultContainerName + strconv.Itoa(containerIndex)
		createContainerAndFiles(containerName)
		fmt.Println(containerName)
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
