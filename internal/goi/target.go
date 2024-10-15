package goi

import (
	"path/filepath"
	"strings"
)

type Target struct {
	FileName string `json:"filename"`
	Arch     string `json:"arch"`
	OS       string `json:"os"`
	Version  string `json:"version"`
}

func (t Target) IsArchive() bool {
	return strings.HasSuffix(t.FileName, ".tar.gz") || filepath.Ext(t.FileName) == ".zip"
}
