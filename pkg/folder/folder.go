package folder

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/idelchi/godyl/pkg/file"
	"github.com/idelchi/godyl/pkg/files"
	"github.com/idelchi/godyl/pkg/utils"
)

// Folder represents a file system directory as a string.
// It provides methods for working with directories, such as creating,
// removing, expanding paths, and checking existence.
type Folder string

// New creates a new Folder object from the provided path segments by joining them.
func New(paths ...string) Folder {
	return Folder(filepath.Join(paths...))
}

// NewInTempDir assigns but does not create a directory inside the system's temporary directory.
func NewInTempDir(paths ...string) Folder {
	return New(os.TempDir(), filepath.Join(paths...))
}

// CreateRandomInDir creates a new random directory inside the given directory.
// Use "" to create a random directory in the default directory for temporary files.
func CreateRandomInDir(dir string, pattern string) (Folder, error) {
	if err := New(dir).CreateIgnoreExisting(); err != nil {
		return Folder(""), fmt.Errorf("creating temporary directory in %s: %w", dir, err)
	}

	name, err := os.MkdirTemp(dir, pattern)
	if err != nil {
		return Folder(""), fmt.Errorf("creating temporary directory in %s: %w", dir, err)
	}

	return Folder(name), nil
}

// CreateIgnoreExisting creates the Folder and all necessary parent directories with 0755 permissions.
// If the directory already exists, the error is nil.
func (f Folder) CreateIgnoreExisting() error {
	const perm = 0o755

	if err := os.MkdirAll(f.Path(), perm); err != nil && !os.IsExist(err) {
		return fmt.Errorf("creating directory %s: %w", f.Path(), err)
	}

	return nil
}

// Create creates the Folder and all necessary parent directories with 0755 permissions.
func (f Folder) Create() error {
	const perm = 0o755

	if err := os.MkdirAll(f.Path(), perm); err != nil {
		return fmt.Errorf("creating directory %s: %w", f.Path(), err)
	}

	return nil
}

// Normalized converts the folder path to use forward slashes.
func (f Folder) Normalized() Folder {
	return New(filepath.ToSlash(f.Path()))
}

// IsSet checks whether the Folder has been set to a non-empty value.
func (f Folder) IsSet() bool {
	return f != ""
}

// IsParentOf determines if the Folder is a parent directory of the given 'other' Folder.
func (f Folder) IsParentOf(other Folder) bool {
	return strings.HasPrefix(other.Path(), f.Path())
}

// Expanded expands the file path in case of ~ and returns the expanded path.
func (f Folder) Expanded() Folder {
	return New(utils.ExpandHome(f.Path()))
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

// Base returns the last element of the Folder's path.
func (f Folder) Base() string {
	return filepath.Base(f.Path())
}

// Remove deletes the Folder and all of its contents.
func (f Folder) Remove() error {
	if err := os.RemoveAll(f.Path()); err != nil {
		return fmt.Errorf("removing directory %s: %w", f.Path(), err)
	}

	return nil
}

// CriteriaFunc defines a function type for filtering files during search operations.
type CriteriaFunc func(file.File) (bool, error)

// ErrNotFound is returned when a file is not found during a search operation.
var ErrNotFound = errors.New("file not found")

// FindFile searches for a file in the Folder that matches the provided criteria.
// It returns the first file found or an error if none are found.
func (f Folder) FindFile(criteria ...CriteriaFunc) (file.File, error) {
	var foundPath file.File

	var found bool

	err := filepath.Walk(f.Path(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if file.File(path).IsDir() {
			return nil // Skip directories
		}

		path, err = filepath.Rel(f.Path(), path)
		if err != nil {
			return err
		}

		// Check if the file matches all criteria
		for _, criterion := range criteria {
			matches, err := criterion(file.File(path))
			if err != nil {
				return err
			}

			if !matches {
				return nil // Skip this file if it doesn't match all criteria
			}
		}

		foundPath = file.New(f.Path(), path)

		// If we've reached here, the file matches all criteria
		found = true

		return filepath.SkipAll // Stop the walk, we've found a match
	})
	if err != nil {
		return file.File(""), fmt.Errorf("error walking folder %q: %w", f.Path(), err)
	}

	if !found {
		return file.File(""), fmt.Errorf("%w: no file found matching all criteria in folder %q", ErrNotFound, f.Path())
	}

	return foundPath, nil
}

// ListFolders returns a slice of Folders representing all subdirectories
// within the current Folder.
func (f Folder) ListFolders() ([]Folder, error) {
	var folders []Folder

	entries, err := os.ReadDir(f.Path())
	if err != nil {
		return nil, fmt.Errorf("reading directory %s: %w", f.Path(), err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			subFolder := New(f.Path(), entry.Name())
			folders = append(folders, subFolder)
		}
	}

	return folders, nil
}

// ListFiles returns a slice of Files representing all files
// within the current Folder.
func (f Folder) ListFiles() (files.Files, error) {
	var files files.Files

	entries, err := os.ReadDir(f.Path())
	if err != nil {
		return nil, fmt.Errorf("reading directory %s: %w", f.Path(), err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			file := file.New(f.Path(), entry.Name())
			files = append(files, file)
		}
	}

	return files, nil
}
