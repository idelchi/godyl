package folder

import (
	"os"
	"path/filepath"

	"github.com/idelchi/godyl/pkg/path/file"
)

// New creates a Folder from one or more path components.
// Joins the paths using filepath.Join and normalizes the result to use forward slashes.
// Note: This does not create the directory, only constructs the path.
func New(paths ...string) Folder {
	return Folder(filepath.ToSlash(filepath.Join(paths...)))
}

// NewInTempDir creates a Folder path in the system temp directory.
// Combines the system temp directory with the provided path components.
// Note: This does not create the directory, only constructs the path.
func NewInTempDir(paths ...string) Folder {
	return New(os.TempDir(), filepath.Join(paths...))
}

// FromFile creates a Folder from a file's parent directory.
// Extracts the directory component from the given file path.
// See `New` for details on path normalization.
// Note: This does not create the directory, only constructs the path.
func FromFile(f file.File) Folder {
	return New(f.Dir())
}
