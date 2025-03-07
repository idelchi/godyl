// Package embed centralizes the management of embedded files used by the application.
package embed

import (
	"embed"
)

// Files holds static template scripts and configuration files.
//
//go:embed ../../defaults.yml
//go:embed ../../tools.yml
//go:embed ../../internal/core/updater/scripts/*
var Files embed.FS

// GetDefaultsFile returns the content of the embedded defaults.yml file.
func GetDefaultsFile() ([]byte, error) {
	return Files.ReadFile("defaults.yml")
}

// GetToolsFile returns the content of the embedded tools.yml file.
func GetToolsFile() ([]byte, error) {
	return Files.ReadFile("tools.yml")
}

// GetScriptFile returns the content of an embedded script file.
func GetScriptFile(name string) ([]byte, error) {
	return Files.ReadFile("internal/core/updater/scripts/" + name)
}
