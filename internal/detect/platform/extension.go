package platform

import (
	"path/filepath"
	"strings"
)

// Extension represents a platform-specific file extension.
// Used primarily for executable files and archive formats.
type Extension string

// Default returns the platform's standard executable extension.
// Returns ".exe" for Windows systems and empty string for Unix-like systems.
func (e *Extension) Default(os OS) Extension {
	switch os.Type {
	case "windows":
		return Extension(".exe")
	default:
		return Extension("")
	}
}

// String returns the extension value including the leading dot.
func (e Extension) String() string {
	return string(e)
}

// Parse extracts the file extension from a filename.
// Handles special cases like ".tar.gz" compound extensions.
func (e *Extension) Parse(name string) error {
	switch ext := filepath.Ext(name); ext {
	case ".gz":
		if strings.HasSuffix(name, ".tar.gz") {
			*e = Extension(".tar.gz")
		} else {
			*e = Extension(ext)
		}
	default:
		*e = Extension(ext)
	}

	return nil
}
