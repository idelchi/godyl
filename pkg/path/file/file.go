package file

import (
	"path/filepath"
)

// New creates a File from one or more path components.
// Joins the paths using filepath.Join and normalizes the result.
func New(paths ...string) File {
	return File(filepath.ToSlash(filepath.Join(paths...)))
}
