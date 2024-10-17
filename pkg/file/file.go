package file

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// File represents a file path as a string, providing methods for file operations.
type File string

// NewFile creates a new File by joining the provided paths.
func NewFile(paths ...string) File {
	return File(filepath.Join(paths...))
}

// Normalize converts the file path to use forward slashes.
func (f *File) Normalize() {
	*f = File(filepath.ToSlash(f.Name()))
}

// Normalize converts the file path to use forward slashes.
func (f File) Normalized() File {
	return File(filepath.ToSlash(f.Name()))
}

// Create creates a new file and returns a pointer to the os.File object, or an error.
func (f File) Create() (*os.File, error) {
	return os.Create(f.String())
}

// Open opens the file for reading and returns a pointer to the os.File object, or an error.
func (f File) Open() (*os.File, error) {
	return os.Open(f.String())
}

// Name returns the name (string representation) of the File.
func (f File) Name() string {
	return f.String()
}

// String returns the string representation of the File.
func (f File) String() string {
	return string(f)
}

// Dir returns the file.Folder object representing the directory of the file.
func (f File) Dir() Folder {
	if f.IsDir() {
		return NewFolder(f.String())
	}

	return NewFolder(filepath.Dir(f.String()))
}

// IsExecutable checks if the file has executable permissions.
func (f File) IsExecutable() (bool, error) {
	info, err := os.Stat(f.String())
	if err != nil {
		return false, fmt.Errorf("getting file info: %w", err)
	}

	return info.Mode()&0o111 != 0, nil
}

// Chmod changes the file permissions to the specified fs.FileMode.
func (f *File) Chmod(mode fs.FileMode) error {
	return os.Chmod(f.String(), mode)
}

// Copy copies the contents of the current file to the specified destination File.
// It returns an error if the operation fails.
func (f File) Copy(other File) error {
	// Open the source file
	source, err := f.Open()
	if err != nil {
		return fmt.Errorf("opening source file: %w", err)
	}
	defer source.Close()

	// Create the destination file
	destination, err := other.Create()
	if err != nil {
		return fmt.Errorf("creating destination file: %w", err)
	}
	defer destination.Close()

	// Copy the contents of the source file to the destination file
	_, err = io.Copy(destination, source)
	if err != nil {
		return fmt.Errorf("copying file: %w", err)
	}

	// Set permissions on the destination file (executable permission)
	if err := destination.Chmod(0o755); err != nil {
		return fmt.Errorf("setting permissions: %w", err)
	}

	return nil
}

// Exists checks if the file exists in the file system.
func (f File) Exists() bool {
	_, err := os.Stat(f.String())
	if err != nil {
		return false // File does not exist or error accessing it
	}

	return true
}

// IsFile checks if the path is a regular file (not a directory or special file).
func (f File) IsFile() bool {
	info, err := os.Stat(f.String())
	if err != nil {
		return false // File does not exist or error accessing it
	}

	return info.Mode().IsRegular()
}

// IsDir checks if the path represents a directory.
func (f File) IsDir() bool {
	info, err := os.Stat(f.String())
	if err != nil {
		return false // File does not exist or error accessing it
	}

	return info.Mode().IsDir()
}

// Extension returns the file extension of the File, mapped to a predefined Extension constant.
func (f File) Extension() Extension {
	ext := filepath.Ext(f.Name())

	switch strings.ToLower(ext) {
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