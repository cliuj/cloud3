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
	clientList = map[string]string{
		"1": "http://localhost:8080",
		"2": "http://localhost:8081",
	}
)

func loadENVs() {
	if SERVER_PORT == "" {
		SERVER_PORT = "8000"
	}

	if SHARED_DIR == "" {
		SHARED_DIR = "/tmp/cloud3/server"
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

	//r.POST("/client/register", func(c *gin.Context) {
	//	id := c.Query("id")
	//	host := c.Request.Host
	//	fmt.Println("host", host)
	//	UpdateClientList(id, host)
	//	fmt.Println(clientList.Load(id))
	//	c.String(http.StatusOK, fmt.Sprintf("Client ID: %s connected from: %s", id, host))
	//})

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

		origin := c.Request.Host

		for _, host := range clientList {
			fmt.Println("trying to send to client: ", origin, host)
			if origin == host {
				// Skip the origin
				continue
			}
			filePaths, err := dirsync.GetFilePathsFromDir(SHARED_DIR)
			if err != nil {
				log.Println(fmt.Errorf("Failed to retrieve filepaths from directory: %s, %v", SHARED_DIR, err))
			}
			uploadURL := host + "/files/upload"
			err = dirsync.UploadFiles(uploadURL, filePaths)
			if err != nil {
				log.Println(fmt.Errorf("Failed to UploadFiles to URL: %s with files %v, %v", uploadURL, filePaths, err))
			}
		}
		c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})

	r.Run(":" + SERVER_PORT)
}
