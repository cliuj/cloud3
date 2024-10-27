package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	CLIENT_PORT = os.Getenv("CLIENT_PORT")
	SHARED_DIR = os.Getenv("SHARED_DIR")
)

func main() {
	fmt.Println("Hello world")

	if CLIENT_PORT == "" {
		CLIENT_PORT = "8080"
	}

	if SHARED_DIR == "" {
		SHARED_DIR = "/tmp/cloud3"
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run(":" + CLIENT_PORT)
}
