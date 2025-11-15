package folder

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/idelchi/godyl/pkg/generic"
	"github.com/idelchi/godyl/pkg/path/file"
)

// WithFile creates a File path within this directory.
// Combines the directory path with the provided filename.
func (f Folder) WithFile(path string) file.File {
	return file.New(f.Path(), path)
}

// IsSet checks if the folder path is non-empty.
func (f Folder) IsSet() bool {
	return f != ""
}

// Join combines this path with additional components.
// Returns a new Folder with the combined path.
func (f Folder) Join(paths ...string) Folder {
	return New(append([]string{f.Path()}, paths...)...)
}

// Expanded resolves environment variables including `~` home directory references.
func (f Folder) Expanded() Folder {
	return New(os.ExpandEnv(generic.ExpandHome(f.Path())))
}

// String returns the folder path as a string.
func (f Folder) String() string {
	return string(f)
}

// Path returns the string representation of the folder path.
func (f Folder) Path() string {
	return f.String()
}

// Exists checks if the directory exists in the filesystem.
func (f Folder) Exists() bool {
	info, err := os.Stat(f.Path())

	return err == nil && info.IsDir()
}

// Base returns the last component of the folder path.
func (f Folder) Base() string {
	return filepath.Base(f.Path())
}

// Up returns the parent directory of this folder.
func (f Folder) Up() Folder {
	return New(filepath.Dir(f.Path()))
}

// RelativeTo returns the path from base to this folder.
// Returns an error if it can't be computed.
func (f Folder) RelativeTo(base Folder) (Folder, error) {
	rel, err := filepath.Rel(base.Path(), f.Path())
	if err != nil {
		return f, fmt.Errorf("getting relative between %q and %q: %w", base, f, err)
	}

	return New(rel), nil
}

// AsFile converts the folder path to a File type.
func (f Folder) AsFile() file.File {
	return file.New(f.Path())
}
