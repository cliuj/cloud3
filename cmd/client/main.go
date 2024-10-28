package main

import (
	"fmt"
	"os"
	"log"
	"path"
	"net/http"
	"io"
	"bytes"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

var (
	CLIENT_PORT = os.Getenv("CLIENT_PORT")
	SERVER_URL = os.Getenv("SERVER_URL")
	SHARED_DIR = os.Getenv("SHARED_DIR")
)

// GetFilePathsFromDir gets the filepaths from a passed directory. The filepath
// is the concatenation of `sourceDir` and the file name.
// NOTE: It will skip processing on directories
//
func GetFilePathsFromDir(sourceDir string) ([]string, error) {
	files, err := os.ReadDir(sourceDir)
	if err != nil {
		log.Println(fmt.Errorf("Error reading directory %s", sourceDir))
		return []string{}, err
	}

	fps := []string{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fp := path.Join(sourceDir, file.Name())
		fps = append(fps, fp)
	}
	return fps, nil
}

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

func main() {
	fmt.Println("Hello world")
	setDefaultEnvs()
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	filePaths, err := GetFilePathsFromDir(SHARED_DIR)
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
