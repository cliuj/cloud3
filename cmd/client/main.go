package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/cliuj/cloud3/internal/dirsync"
)

var (
	CLIENT_PORT = os.Getenv("CLIENT_PORT")
	SERVER_URL = os.Getenv("SERVER_URL")
	SHARED_DIR = os.Getenv("SHARED_DIR")
)

// UploadFiles attempts to upload the files specified in filePaths via a POST request
//
func UploadFiles(destURL string, filePaths []string) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for _, filePath := range filePaths {
		f, err := os.Open(filePath)
		if err != nil {
			log.Println(fmt.Sprintf("Failed to open file: %s", filePath))
			return err
		}
		defer f.Close()

		part, err := writer.CreateFormFile("upload[]", filePath)
		if err != nil {
			log.Println(fmt.Sprintf("Failed to CreateFormFile for: %s", filePath))
			return err
		}
		_, err = io.Copy(part, f)
		if err != nil {
			log.Println(fmt.Sprintf("Failed to copy file: %s", filePath))
			return err
		}

	}
	writer.Close()
	resp, err := http.Post(destURL, writer.FormDataContentType(), body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil

}

func setDefaultEnvs() {
	if CLIENT_PORT == "" {
		CLIENT_PORT = "8080"
	}

	if SHARED_DIR == "" {
		SHARED_DIR = "/tmp/cloud3/client"
	}

	if SERVER_URL == "" {
		SERVER_URL = "http://localhost:8000"
	}
}

func PollDirChanges(sourceDir string) {
	for {
		checksum, err := dirsync.GetDirSHASUM(sourceDir)
		if err != nil {
			// TODO: Need to handle this later
			log.Println(fmt.Errorf("Error while retrieving checksum of dir %s, %v", sourceDir, err))
		}
		fmt.Println(checksum)

		// if checksum diff:
		// upload to server

		time.Sleep(time.Second)
	}
}

func main() {
	fmt.Println("Hello world")
	setDefaultEnvs()
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})


	//go GetDirSHASUM(SHARED_DIR)
	go PollDirChanges(SHARED_DIR)


	filePaths, err := dirsync.GetFilePathsFromDir(SHARED_DIR)
	if err != nil {
		log.Fatalf("Failed to retrieve filepaths from directory: %s, %v", SHARED_DIR, err)
	}
	uploadURL := SERVER_URL + "/files/upload"
	err = UploadFiles(uploadURL, filePaths)
	if err != nil {
		log.Fatalf("Failed to UploadFiles to URL: %s with files %v, %v", uploadURL, filePaths, err)
	}

	r.Run(":" + CLIENT_PORT)
}
