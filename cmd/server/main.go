package main

import (
	"fmt"
	"os"
	"github.com/gin-gonic/gin"
)

var (
	SERVER_PORT = os.Getenv("SERVER_PORT")
)


func main() {
	fmt.Println("Hello world")
	if SERVER_PORT == "" {
		SERVER_PORT = "8000"
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(":" + SERVER_PORT)
}
