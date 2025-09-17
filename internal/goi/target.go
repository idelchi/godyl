package goi

import (
	"path/filepath"
	"strings"
)

// Target represents a downloadable file in a Go release, including its filename, architecture, OS, and version.
type Target struct {
	FileName string `json:"filename"` // FileName is the name of the file associated with the target.
	Arch     string `json:"arch"`     // Arch specifies the architecture (e.g., amd64, arm) the file is built for.
	OS       string `json:"os"`       // OS specifies the operating system (e.g., linux, windows) the file is intended for.
	Version  string `json:"version"`  // Version is the version of the Go release associated with this file.
	Checksum string `json:"sha256"`   // Checksum is the SHA256 checksum of the file for integrity verification.
}

// IsArchive checks if the target file is an archive, either a .tar.gz or .zip file.
func (t Target) IsArchive() bool {
	return strings.HasSuffix(t.FileName, ".tar.gz") || filepath.Ext(t.FileName) == ".zip"
}
