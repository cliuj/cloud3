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
	"encoding/json"

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
		// Get local checksum
		clientDirChecksum, err := dirsync.GetDirSHASUM(sourceDir)
		if err != nil {
			// TODO: Need to handle this later
			log.Println(fmt.Errorf("Error while retrieving checksum of dir %s, %v", sourceDir, err))
		}
		fmt.Println(clientDirChecksum)

		// Get Remote checksum
		requestURL := SERVER_URL + "/checksum"
		serverChecksum, err := GetServerChecksum(requestURL)
		fmt.Println(serverChecksum)

		log.Println("client:", clientDirChecksum, "server:", serverChecksum)

		// If both checksums do not match, then upload local to remote
		if clientDirChecksum != serverChecksum {
			filePaths, err := dirsync.GetFilePathsFromDir(SHARED_DIR)
			if err != nil {
				log.Fatalf("Failed to retrieve filepaths from directory: %s, %v", SHARED_DIR, err)
			}
			uploadURL := SERVER_URL + "/files/upload"
			err = UploadFiles(uploadURL, filePaths)
			if err != nil {
				log.Fatalf("Failed to UploadFiles to URL: %s with files %v, %v", uploadURL, filePaths, err)
			}
		}
		time.Sleep(time.Second)
	}
}

func GetServerChecksum(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		// TODO: log message
		return "", fmt.Errorf("Server checksum returned Non-200 code: %v", resp.StatusCode)
	}
	defer resp.Body.Close()
	data := dirsync.ChecksumPayload{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data.Checksum, nil
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

	go PollDirChanges(SHARED_DIR)

	r.Run(":" + CLIENT_PORT)
}
