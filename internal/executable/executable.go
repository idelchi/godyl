package executable

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

// Executable consists of a full path to a file and its version.
// An attempt to parse the version into a string can be done by using a `Version` type.
type Executable struct {
	Path    string
	Version string
}

// New creates a new Executable with the given directory and path.
// The path is joined with the directory to create the full path.
func New(dir string, path string) Executable {
	return Executable{Path: filepath.Join(dir, path)}
}

// Name returns the base name of the executable.
func (e Executable) Name() string {
	return filepath.Base(e.Path)
}

// Dir returns the directory of the executable.
func (e Executable) Dir() string {
	return filepath.Dir(e.Path)
}

// CriteriaFunc is a function that takes an Executable and returns a boolean and an error.
type CriteriaFunc func(Executable) (bool, error)

// Find searches for a file in the given directory that matches all provided criteria.
func (e Executable) Find(dir string, criteria ...CriteriaFunc) (Executable, error) {
	var filePath string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		match := regexp.MustCompile(e.Name()).MatchString(info.Name())

		if match && !info.IsDir() {
			if len(criteria) == 0 {
				// If no criteria provided, accept the first file matching the name
				filePath = path
				return filepath.SkipDir
			}

			for _, criterion := range criteria {
				matches, err := criterion(e)
				if err != nil {
					return err
				}
				if !matches {
					return nil // Skip this file if it doesn't match all criteria
				}
			}
			filePath = path
			return filepath.SkipDir // Stop walking, we found a match
		}
		return nil
	})
	if err != nil {
		return e, err
	}

	if filePath == "" {
		return e, fmt.Errorf("file %q not found matching all criteria", e.Name())
	}
	return Executable{Path: filePath}, nil
}

// IsExecutable checks if the file is executable.
func IsExecutable(file Executable) (bool, error) {
	info, err := os.Stat(file.Path)
	if err != nil {
		return false, fmt.Errorf("getting file info: %w", err)
	}

	return info.Mode()&0o111 != 0, nil
}

// Copy copies the file to the destination path and sets the executable permission.
func (e Executable) Copy(destination string) error {
	// Open the source file
	sourceFile, err := os.Open(e.Path)
	if err != nil {
		return fmt.Errorf("opening source file: %w", err)
	}
	defer sourceFile.Close()

	// Create the destination file
	destinationFile, err := os.Create(destination)
	if err != nil {
		return fmt.Errorf("creating destination file: %w", err)
	}
	defer destinationFile.Close()

	// Copy the contents of the source file to the destination file
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("copying file: %w", err)
	}

	// Set permissions on the destination file (executable permission)
	err = os.Chmod(destination, 0o755)
	if err != nil {
		return fmt.Errorf("setting permissions: %w", err)
	}

	return nil
}

// Exists checks if the file exists.
func (e Executable) Exists() bool {
	info, err := os.Stat(e.Path)
	if err != nil {
		return false // File does not exist or error accessing it
	}

	// Check if the path is a regular file (not a folder or special file)
	return info.Mode().IsRegular()
}
