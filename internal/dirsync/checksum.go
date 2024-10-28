package dirsync

import (
	"os"
	"log"
	"fmt"
	"path"
	"io"
	"crypto/sha256"
	"encoding/hex"
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

func GetFileSHASUM(filePath string) ([]byte, error) {

	f, err := os.Open(filePath)
	if err != nil {
		return []byte{}, err
	}

	defer f.Close()
	h := sha256.New()

	_, err = io.Copy(h, f)
	if err != nil {
		return []byte{}, err
	}

	//log.Println("h:", h)

	return h.Sum(nil), nil
}

func GetDirSHASUM(sourceDir string) (string, error) {
	filePaths, err := GetFilePathsFromDir(sourceDir)
	if err != nil {
		return "", err
	}
	
	checksum := sha256.New()

	for _, filePath := range filePaths {
		sum, err := GetFileSHASUM(filePath)
		if err != nil {
			return "", err
		}
		checksum.Write(sum)
	}

	dirChecksum := hex.EncodeToString(checksum.Sum(nil))
	log.Println("dirChecksum:", dirChecksum)
	return dirChecksum, nil
}
