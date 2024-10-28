package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"time"
	"bytes"

	"github.com/gin-gonic/gin"

	"github.com/cliuj/cloud3/internal/client"
	"github.com/cliuj/cloud3/internal/dirsync"
)

var (
	CLIENT_PORT = os.Getenv("CLIENT_PORT")
	CLIENT_ID = os.Getenv("CLIENT_ID")
	SERVER_URL = os.Getenv("SERVER_URL")
	SHARED_DIR = os.Getenv("SHARED_DIR")
)

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

	if CLIENT_ID == "" {
		log.Fatalf("CLIENT_ID needs to be set!")
	}
}

// PollDirChanges retrieves the checksum for both the locally passed directory
// as well as the remote server's checksum, compares the 2 checksums and attempts to
// upload the local state of the directory to the remote server.
//
func PollDirChanges(sourceDir string) {
	for {
		// Get local checksum
		clientDirChecksum, err := dirsync.GetDirSHASUM(sourceDir)
		if err != nil {
			// TODO: Need to handle this later
			log.Println(fmt.Errorf("Error while retrieving checksum of dir %s, %v", sourceDir, err))
		}

		// Get Remote checksum
		requestURL := SERVER_URL + "/checksum"
		serverChecksum, err := GetServerChecksum(requestURL)

		// If both checksums do not match, then upload local to remote
		if clientDirChecksum != serverChecksum {
			filePaths, err := dirsync.GetFilePathsFromDir(SHARED_DIR)
			if err != nil {
				log.Fatalf("Failed to retrieve filepaths from directory: %s, %v", SHARED_DIR, err)
			}
			uploadURL := SERVER_URL + "/files/upload"
			err = dirsync.UploadFiles(uploadURL, filePaths)
			if err != nil {
				log.Fatalf("Failed to UploadFiles to URL: %s with files %v, %v", uploadURL, filePaths, err)
			}
		}
		time.Sleep(time.Second)
	}
}

// GetServerChecksum is wrapper function to make a simple HTTP request call
// to get the remote server's checksum.
//
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

// RegisterClient is a wrapper function to make a simple HTTP request call to the
// remote server so that the server can note down which clients it needs to connect
// and run upload on.
func RegisterClient(url string) error {
	jsonPayload, err := json.Marshal(
		client.RegisterClientPayload{
		ID: CLIENT_ID,
	})
	if err != nil {
		log.Println("Error unmarshalling JSON", err)
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Println("Error registering client", err)
		return err
	}
	defer resp.Body.Close()
	return nil
}



func main() {
	setDefaultEnvs()
	//RegisterClient(SERVER_URL + "/client/register")
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	go PollDirChanges(SHARED_DIR)

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

	r.Run(":" + CLIENT_PORT)
}
