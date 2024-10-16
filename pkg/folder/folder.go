package folder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Folder represents a file system directory as a string.
// It provides methods for working with directories, such as creating,
// removing, expanding paths, and checking existence.
type Folder string

// New creates a new Folder from the provided path segments by joining them.
func New(paths ...string) Folder {
	return Folder(filepath.Join(paths...))
}

// IsSet checks whether the Folder has been set to a non-empty value.
func (f Folder) IsSet() bool {
	return f != ""
}

// Set assigns a new path to the Folder.
func (f *Folder) Set(path string) {
	*f = Folder(path)
}

// IsParentOf determines if the Folder is a parent directory of the given 'other' Folder.
func (f *Folder) IsParentOf(other Folder) bool {
	return strings.HasPrefix(other.Path(), f.Path())
}

// Expand expands a Folder path that begins with "~" to the user's home directory.
func (f *Folder) Expand() error {
	if strings.HasPrefix(f.Path(), "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("getting user home directory: %w", err)
		}
		f.Set(filepath.Join(homeDir, f.Path()[2:]))
	}
	return nil
}

// String returns the Folder as a string.
func (f Folder) String() string {
	return string(f)
}

// Path returns the Folder path as a string.
func (f Folder) Path() string {
	return f.String()
}

// Exists checks if the Folder exists in the file system.
func (f Folder) Exists() bool {
	_, err := os.Stat(f.Path())
	return err == nil
}

// Create creates the Folder and all necessary parent directories with 0755 permissions.
func (f Folder) Create() error {
	return os.MkdirAll(f.Path(), 0o755)
}

// Name returns the base name (last element) of the Folder's path.
func (f Folder) Name() string {
	return filepath.Base(f.Path())
}

// CreateRandomInTempDir creates a new random directory inside the system's temporary directory
// and assigns the generated path to the Folder.
func (f *Folder) CreateRandomInTempDir() error {
	name, err := os.MkdirTemp("", "godyl-*")
	f.Set(name)
	return err
}

// CreateInTempDir creates a directory inside the system's temporary directory
// using the Folder's name and assigns the path to the Folder.
func (f *Folder) CreateInTempDir() error {
	name := filepath.Join(os.TempDir(), f.Name())
	err := os.Mkdir(name, 0o755)
	f.Set(name)
	return err
}

// Remove deletes the Folder and all of its contents.
func (f Folder) Remove() error {
	return os.RemoveAll(f.Path())
}
