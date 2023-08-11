package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func coderunnerProcesser(c *gin.Context) {
	var req struct {
		Code string `json:"code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "无法解析 JSON 数据"})
		return
	}

	fmt.Println(req.Code)
	Response, Find := findCache(req.Code)
	if Find {
		//fmt.Println("i have found cache!")
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
