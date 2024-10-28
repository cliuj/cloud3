package main

import (
	"fmt"
	"os"
	"log"
	"path"
	"net/http"
	"github.com/gin-gonic/gin"
)

var (
	SERVER_PORT = os.Getenv("SERVER_PORT")
	SHARED_DIR = os.Getenv("SHARED_DIR")
)

func loadENVs() {
	if SERVER_PORT == "" {
		SERVER_PORT = "8080"
	}

	if SHARED_DIR == "" {
		SHARED_DIR = "/tmp/cloud3/client"
	}

}

func main() {
	fmt.Println("Hello world")
	loadENVs()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	
	r.POST("/files/upload", func(c *gin.Context) {
		form, _ := c.MultipartForm()
		files := form.File["upload[]"]

		for _, file := range files {
			log.Println(file.Filename)

			dst := path.Join(SHARED_DIR, file.Filename)
			c.SaveUploadedFile(file, dst)
		}
		c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})

	r.Run(":" + SERVER_PORT)
}
