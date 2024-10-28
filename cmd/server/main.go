package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/cliuj/cloud3/internal/dirsync"
	"github.com/gin-gonic/gin"
)

var (
	SERVER_PORT = os.Getenv("SERVER_PORT")
	SHARED_DIR = os.Getenv("SHARED_DIR")
)

func loadENVs() {
	if SERVER_PORT == "" {
		SERVER_PORT = "8000"
	}

	if SHARED_DIR == "" {
		SHARED_DIR = "/tmp/cloud3/server"
	}

}

func getChecksum(sourceDir string) {
	
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

	r.GET("/checksum", func(c *gin.Context) {
		checksum, err := dirsync.GetDirSHASUM(SHARED_DIR)
		if err != nil {
			log.Println("Failed to get checksum", err)
			c.JSON(500, gin.H{
				"message": "failed to get checksum",
			})
			return
		}
		c.JSON(200, gin.H{
			"message": checksum,
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
