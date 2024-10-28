package dirsync

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"log"
	"fmt"
	"net/http"
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
