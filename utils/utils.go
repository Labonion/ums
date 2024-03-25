package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, result interface{}, statusCode int){
	c.JSON(statusCode, gin.H{
		"status": statusCode,
		"data" : result,
		"success": true,
	})
}

func Failure(c *gin.Context, err  interface{}, statusCode int64){
	fmt.Println(err)
	c.JSON(int(statusCode), gin.H{
		"status": statusCode,
		"error" : err,
		"success": false,
	})
}