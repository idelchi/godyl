package match

import (
	"path/filepath"
)

// Extension represents a file extension.
type Extension int

const (
	// None represents no file extension.
	None Extension = iota
	// EXE represents the ".exe" file extension.
	EXE
	// GZ represents the ".gz" file extension.
	GZ
	// ZIP represents the ".zip" file extension.
	ZIP
	// Other represents any other file extension.
	Other
)

// Extension returns the file extension of the asset based on its name.
// It maps common file extensions to predefined constants.
func (a *Asset) Extension() Extension {
	ext := filepath.Ext(a.NameLower())

	switch ext {
	case ".exe":
		return EXE
	case ".gz":
		return GZ
	case ".zip":
		return ZIP
	case "":
		return None
	default:
		return Other
	}
}
