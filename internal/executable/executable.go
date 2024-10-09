package executable

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Executable consists of a full path to a file and its version-
type Executable struct {
	Path    string
	Version string
}

func New(dir string, path string) Executable {
	return Executable{Path: filepath.Join(dir, path)}
}

func (e Executable) Name() string {
	return filepath.Base(e.Path)
}

func (e Executable) Dir() string {
	return filepath.Dir(e.Path)
}

type CriteriaFunc func(Executable) (bool, error)

// Find searches for a file in the given directory that matches all provided criteria
func (e Executable) Find(dir string, criteria ...CriteriaFunc) (Executable, error) {
	var filePath string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == e.Name() && !info.IsDir() {
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

func IsExecutable(file Executable) (bool, error) {
	info, err := os.Stat(file.Path)
	if err != nil {
		return false, fmt.Errorf("getting file info: %w", err)
	}

	return info.Mode()&0o111 != 0, nil
}

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

func (e Executable) Exists() bool {
	info, err := os.Stat(e.Path)
	if err != nil {
		return false // File does not exist or error accessing it
	}

	// Check if the path is a regular file (not a folder or special file)
	return info.Mode().IsRegular()
}
