package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/goccy/go-json"
)

type FileStruct struct {
	ProjectID int  `json:"projectID"`
	FileID    int  `json:"fileID"`
	Required  bool `json:"required"`
}
type ManifestStruct struct {
	Files []FileStruct `json:"files"`
	Name  string       `json:"name"`
}

func ERROR_INVALID_MANIFEST_FILE() { log.Println("Invalid Manifest File") }

func CFDownloadURLBuilder(v FileStruct) string {
	return fmt.Sprintf("https://www.curseforge.com/api/v1/mods/%d/files/%d/download", v.ProjectID, v.FileID)
}

func DownloadFiles(files []FileStruct, i int, sleepTime uint) {
	if i+1 == len(files) {
		fmt.Printf("Finished running - %d/%d files downloaded", i+1, len(files))
		return
	}
	v := files[i]
	url := CFDownloadURLBuilder(v)
	outputFile, err := os.Create(fmt.Sprintf("./mods/%d.jar", v.FileID))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Downloading file %d/%d\n", i+1, len(files))
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	_, err = io.Copy(outputFile, resp.Body)
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Duration(sleepTime * 1000))
	DownloadFiles(files, i+1, sleepTime)
}

func main() {
	manifestFlag := flag.String("manifest", "", "manifest.json file")
	sleepTimeFlag := flag.Uint("delay", 200, "time between each iteration")
	flag.Parse()

	manifestFile, err := os.ReadFile(*manifestFlag)
	if err != nil {
		ERROR_INVALID_MANIFEST_FILE()
		return
	}

	var manifestContent ManifestStruct
	err = json.Unmarshal(manifestFile, &manifestContent)
	if err != nil {
		ERROR_INVALID_MANIFEST_FILE()
		return
	}

	DownloadFiles(manifestContent.Files, 0, *sleepTimeFlag)
}
