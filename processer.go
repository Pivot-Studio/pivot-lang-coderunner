package main

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func coderunnerProcesser(c *gin.Context) {
	defer func() {
		semaphore <- struct{}{} // 处理完毕后释放信号量
	}()

	<-semaphore

	var req struct {
		Code string `json:"code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无法解析 JSON 数据"})
		return
	}



	containerIndex = containerIndex % containerNum
	containerIndex++

	//containername是default和index粘起来
	containerName = defaultContainerName + strconv.Itoa(containerIndex)
	createContainerAndFiles()
	fmt.Println(containerName)

	fmt.Println(req.Code)
	Response, Found := findCache(req.Code)
	if Found {
		c.JSON(200, Response)
	} else {
		var err error
		Response, err = coderunner(req.Code)
		insertCache(req.Code, Response)
		if err != nil {
			Response.Status = 500
		} else {
			Response.Status = 200
		}
		c.JSON(200, Response)
	}
}
