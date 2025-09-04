package file

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Set checks if the path is non-empty.
func (f File) Set() bool {
	return f.String() != ""
}

// IsExecutable checks if the file has execute permissions.
// Returns true if any execute bit (user/group/other) is set.
func (f File) IsExecutable() (bool, error) {
	info, err := os.Stat(f.String())
	if err != nil {
		return false, fmt.Errorf("getting file info: %w", err)
	}

	return info.Mode()&0o111 != 0, nil
}

// Exists checks if the path exists in the filesystem.
// Returns true if the path exists, false otherwise.
func (f File) Exists() bool {
	_, err := os.Stat(f.String())

	return err == nil
}

// IsFile checks if the path is a regular file.
// Returns false for directories, symlinks, and special files.
func (f File) IsFile() bool {
	info, err := os.Stat(f.String())
	if err != nil {
		return false // File does not exist or error accessing it
	}

	return info.Mode().IsRegular()
}

// IsDir checks if the path is a directory.
// Returns false for regular files and non-existent paths.
func (f File) IsDir() bool {
	info, err := os.Stat(f.String())
	if err != nil {
		return false // File does not exist or error accessing it
	}

	return info.Mode().IsDir()
}

// Extension returns the file's extension as a string, without the leading dot.
func (f File) Extension() string {
	return strings.TrimPrefix(filepath.Ext(f.String()), ".")
}

// Info retrieves the file information.
func (f File) Info() (fs.FileInfo, error) {
	info, err := os.Stat(f.String())
	if err != nil {
		return nil, fmt.Errorf("getting file info for %q: %w", f, err)
	}

	return info, nil
}

// Size returns the size of the file in bytes.
func (f File) Size() (int64, error) {
	info, err := f.Info()
	if err != nil {
		return 0, fmt.Errorf("getting file size for %q: %w", f, err)
	}

	if !f.IsFile() {
		return 0, fmt.Errorf("file %q is not a regular file", f)
	}

	return info.Size(), nil
}

// LargerThan checks if the file is larger than the specified size in bytes.
func (f File) LargerThan(size int64) (bool, error) {
	currentSize, err := f.Size()
	if err != nil {
		return false, fmt.Errorf("checking if file %q is larger than %d bytes: %w", f, size, err)
	}

	return currentSize > size, nil
}

// SmallerThan checks if the file is smaller than the specified size in bytes.
func (f File) SmallerThan(size int64) (bool, error) {
	currentSize, err := f.Size()
	if err != nil {
		return false, fmt.Errorf("checking if file %q is smaller than %d bytes: %w", f, size, err)
	}

	return currentSize < size, nil
}

// Hash computes the hash of the file's contents.
func (f File) Hash() (string, error) {
	data, err := f.Read()
	if err != nil {
		return "", fmt.Errorf("reading file %q: %w", f, err)
	}

	hash := sha256.Sum256(data)

	return hex.EncodeToString(hash[:]), nil
}

// InPath checks if the file can be found in the system PATH.
// Returns true if the file is found in PATH, false otherwise.
func (f File) InPath() bool {
	_, err := exec.LookPath(f.Path())

	return err == nil
}

// NumberOfLines returns the number of lines in the file.
func (f File) NumberOfLines() (int, error) {
	data, err := f.Read()
	if err != nil {
		return 0, fmt.Errorf("reading file %q: %w", f, err)
	}

	return len(strings.Split(string(data), "\n")), nil
}
