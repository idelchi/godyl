package main

import (
	"fmt"
	"log"

	"github.com/idelchi/godyl/internal/folder"
	"github.com/idelchi/godyl/pkg/download"
)

func main() {
	// Create a new Downloader
	d := download.New()

	tmp := folder.Folder("")
	tmp.CreateRandomInTempDir()

	// Example 1: Download a regular file
	fileURL := "https://github.com/jqlang/jq/releases/download/jq-1.7.1/jq-linux-amd64"
	fileDst, err := d.Download(fileURL, tmp.Path())
	if err != nil {
		log.Fatalf("Failed to download file: %v", err)
	}
	fmt.Printf("File downloaded to: %s\n", fileDst)

	tmp.CreateRandomInTempDir()

	// Example 2: Download a zip archive and extract it
	zipURL := "https://github.com/starship/starship/releases/download/v1.20.1/starship-aarch64-apple-darwin.tar.gz"
	zipDst, err := d.Download(zipURL, tmp.Path())
	if err != nil {
		log.Fatalf("Failed to download and extract zip: %v", err)
	}
	fmt.Printf("Zip extracted to: %s\n", zipDst)
}
