package file

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

type File string

func (f File) Name() string {
	return f.String()
}

func (f File) String() string {
	return string(f)
}

func (f File) Dir() string {
	return filepath.Dir(f.String())
}

type CriteriaFunc func(File) (bool, error)

// Find searches for a file in the given directory that matches all provided criteria.
func (f File) Find(dir string, criteria ...CriteriaFunc) (File, error) {
	var filePath string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		match := regexp.MustCompile(f.Name()).MatchString(info.Name())

		if match && !info.IsDir() {
			if len(criteria) == 0 {
				// If no criteria provided, accept the first file matching the name
				filePath = path
				return filepath.SkipDir
			}

			for _, criterion := range criteria {
				matches, err := criterion(f)
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
		return f, err
	}

	if filePath == "" {
		return f, fmt.Errorf("file %q not found matching all criteria", f.Name())
	}
	return File(filePath), nil
}

func (f File) IsExecutable() (bool, error) {
	info, err := os.Stat(f.String())
	if err != nil {
		return false, fmt.Errorf("getting file info: %w", err)
	}

	return info.Mode()&0o111 != 0, nil
}

func (f File) Copy(destination string) error {
	// Open the source file
	sourceFile, err := os.Open(f.String())
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
func (f File) Exists() bool {
	_, err := os.Stat(f.String())
	if err != nil {
		return false // File does not exist or error accessing it
	}

	return true
}

// Exists checks if the file exists.
func (f File) IsFile() bool {
	info, err := os.Stat(f.String())
	if err != nil {
		return false // File does not exist or error accessing it
	}

	// Check if the path is a regular file (not a folder or special file)
	return info.Mode().IsRegular()
}

func (f File) IsDir() bool {
	info, err := os.Stat(f.String())
	if err != nil {
		return false // File does not exist or error accessing it
	}

	// Check if the path is a regular file (not a folder or special file)
	return info.Mode().IsDir()
}
