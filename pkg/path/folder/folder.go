package folder

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/files"
	"github.com/idelchi/godyl/pkg/utils"
)

// Folder represents a file system directory as a string.
// It provides methods for working with directories, such as creating,
// removing, expanding paths, and checking existence.
type Folder string

// New creates a new Folder object from the provided path segments by joining them.
func New(paths ...string) Folder {
	return Folder(filepath.Clean(filepath.Join(paths...))).Normalized()
}

// NewInTempDir assigns but does not create a directory inside the system's temporary directory.
func NewInTempDir(paths ...string) Folder {
	return New(os.TempDir(), filepath.Join(paths...))
}

// FromFile creates a folder from a file's directory path.
func FromFile(f file.File) Folder {
	return New(f.Dir())
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

	return New(name), nil
}

// CreateIgnoreExisting creates the Folder and all necessary parent directories with 0755 permissions.
// If the directory already exists, the error is nil.
func (f Folder) CreateIgnoreExisting() error {
	const perm = 0o755

	if err := os.MkdirAll(f.String(), perm); err != nil && !os.IsExist(err) {
		return fmt.Errorf("creating directory %s: %w", f.String(), err)
	}

	return nil
}

// WithFile returns a new Folder with the provided file name appended to the current folder.
func (f Folder) WithFile(path string) file.File {
	return file.New(f.String(), path)
}

// Create creates the Folder and all necessary parent directories with 0755 permissions.
func (f Folder) Create() error {
	const perm = 0o755

	if err := os.MkdirAll(f.String(), perm); err != nil {
		return fmt.Errorf("creating directory %s: %w", f.String(), err)
	}

	return nil
}

// Normalized converts the folder path to use forward slashes.
func (f Folder) Normalized() Folder {
	return Folder(filepath.ToSlash(f.String()))
}

// IsSet checks whether the Folder has been set to a non-empty value.
func (f Folder) IsSet() bool {
	return f != ""
}

// IsParentOf determines if the Folder is a parent directory of the given 'other' Folder.
func (f Folder) IsParentOf(other Folder) bool {
	return strings.HasPrefix(other.String(), f.String())
}

// Join joins the Folder with the provided paths and returns a new Folder.
func (f Folder) Join(paths ...string) Folder {
	return New(append([]string{f.String()}, paths...)...)
}

// Expanded expands the file path in case of ~ and returns the expanded path.
func (f Folder) Expanded() Folder {
	return New(utils.ExpandHome(f.String()))
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
	_, err := os.Stat(f.String())

	return err == nil
}

// Base returns the last element of the Folder's path.
func (f Folder) Base() string {
	return filepath.Base(f.String())
}

func (f Folder) RelativeTo(base Folder) (Folder, error) {
	file, err := filepath.Rel(f.String(), base.String())
	if err != nil {
		return f, fmt.Errorf("getting relative path: %w", err)
	}

	return New(file), nil
}

// Remove deletes the Folder and all of its contents.
func (f Folder) Remove() error {
	if err := os.RemoveAll(f.String()); err != nil {
		return fmt.Errorf("removing directory %s: %w", f.String(), err)
	}

	return nil
}

func (f Folder) AsFile() file.File {
	return file.New(f.String())
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

	root := f.AsFile()

	err := filepath.Walk(root.String(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		current := file.New(path)

		if current.IsDir() {
			return nil // Skip directories
		}

		relPath, err := root.RelativeTo(current)
		if err != nil {
			return err
		}

		// Check if the file matches all criteria
		for _, criterion := range criteria {
			matches, err := criterion(relPath)
			if err != nil {
				return err
			}

			if !matches {
				return nil // Skip this file if it doesn't match all criteria
			}
		}

		foundPath = root.Join(relPath.String())

		// If we've reached here, the file matches all criteria
		found = true

		return filepath.SkipAll // Stop the walk, we've found a match
	})
	if err != nil {
		return file.File(""), fmt.Errorf("error walking folder %q: %w", root, err)
	}

	if !found {
		return file.File(""), fmt.Errorf("%w: no file found matching all criteria in folder %q", ErrNotFound, root)
	}

	return foundPath, nil
}

// ListFolders returns a slice of Folders representing all subdirectories
// within the current Folder.
func (f Folder) ListFolders() ([]Folder, error) {
	var folders []Folder

	entries, err := os.ReadDir(f.String())
	if err != nil {
		return nil, fmt.Errorf("reading directory %s: %w", f.String(), err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			subFolder := New(f.String(), entry.Name())
			folders = append(folders, subFolder)
		}
	}

	return folders, nil
}

// ListFiles returns a slice of Files representing all files
// within the current Folder.
func (f Folder) ListFiles() (files.Files, error) {
	var files files.Files

	entries, err := os.ReadDir(f.String())
	if err != nil {
		return nil, fmt.Errorf("reading directory %s: %w", f.String(), err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			file := file.New(f.String(), entry.Name())
			files = append(files, file)
		}
	}

	return files, nil
}
