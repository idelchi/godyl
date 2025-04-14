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

// Folder represents a filesystem directory path.
// Provides methods for directory operations including creation,
// removal, path manipulation, and file searching. Handles path
// normalization and expansion of special characters like ~.
type Folder string

// New creates a Folder from one or more path components.
// Joins the paths using filepath.Join and normalizes the result.
func New(paths ...string) Folder {
	return Folder(filepath.Clean(filepath.Join(paths...))).Normalized()
}

// NewInTempDir creates a Folder path in the system temp directory.
// Combines the system temp directory with the provided path components.
// Note: This does not create the directory, only constructs the path.
func NewInTempDir(paths ...string) Folder {
	return New(os.TempDir(), filepath.Join(paths...))
}

// FromFile creates a Folder from a file's parent directory.
// Extracts the directory component from the given file path.
func FromFile(f file.File) Folder {
	return New(f.Dir())
}

// CreateRandomInDir creates a uniquely named directory.
// Creates a directory with a random name inside the specified parent.
// Use empty string for parent to create in system temp directory.
// Pattern is used as a prefix for the random directory name.
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

// CreateIgnoreExisting ensures the directory exists.
// Creates the directory and any missing parent directories.
// Returns nil if the directory already exists.
func (f Folder) CreateIgnoreExisting() error {
	const perm = 0o755

	if err := os.MkdirAll(f.String(), perm); err != nil && !os.IsExist(err) {
		return fmt.Errorf("creating directory %s: %w", f.String(), err)
	}

	return nil
}

// WithFile creates a File path within this directory.
// Combines the directory path with the provided filename.
func (f Folder) WithFile(path string) file.File {
	return file.New(f.String(), path)
}

// Create ensures the directory and its parents exist.
// Creates all necessary directories with 0755 permissions.
func (f Folder) Create() error {
	const perm = 0o755

	if err := os.MkdirAll(f.String(), perm); err != nil {
		return fmt.Errorf("creating directory %s: %w", f.String(), err)
	}

	return nil
}

// Normalized returns the path with forward slashes.
// Converts backslashes to forward slashes for consistency.
func (f Folder) Normalized() Folder {
	return Folder(filepath.ToSlash(f.String()))
}

// IsSet checks if the folder path is non-empty.
func (f Folder) IsSet() bool {
	return f != ""
}

// IsParentOf checks if this folder contains another folder.
// Returns true if the other folder's path starts with this folder's path.
func (f Folder) IsParentOf(other Folder) bool {
	return strings.HasPrefix(other.String(), f.String())
}

// Join combines this path with additional components.
// Returns a new Folder with the combined path.
func (f Folder) Join(paths ...string) Folder {
	return New(append([]string{f.String()}, paths...)...)
}

// Expanded resolves home directory references.
// Replaces ~ with the user's home directory path.
func (f Folder) Expanded() Folder {
	return New(utils.ExpandHome(f.String()))
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
	_, err := os.Stat(f.String())

	return err == nil
}

// Base returns the last component of the folder path.
func (f Folder) Base() string {
	return filepath.Base(f.String())
}

// RelativeTo computes the relative path from this folder to a base path.
// Returns an error if the relative path cannot be computed.
func (f Folder) RelativeTo(base Folder) (Folder, error) {
	file, err := filepath.Rel(f.String(), base.String())
	if err != nil {
		return f, fmt.Errorf("getting relative path: %w", err)
	}

	return New(file), nil
}

// Remove recursively deletes the directory and its contents.
func (f Folder) Remove() error {
	if err := os.RemoveAll(f.String()); err != nil {
		return fmt.Errorf("removing directory %s: %w", f.String(), err)
	}

	return nil
}

// AsFile converts the folder path to a File type.
func (f Folder) AsFile() file.File {
	return file.New(f.String())
}

// CriteriaFunc defines a file matching predicate.
// Returns true if a file matches the criteria, false otherwise.
type CriteriaFunc func(file.File) (bool, error)

// ErrNotFound indicates a file matching the search criteria was not found.
var ErrNotFound = errors.New("file not found")

// FindFile searches for a file matching all given criteria.
// Recursively searches the directory tree and returns the first
// matching file. Returns ErrNotFound if no file matches all criteria.
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

// ListFolders returns all immediate subdirectories.
// Returns a slice of Folders for each directory entry,
// excluding regular files and other filesystem objects.
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

// ListFiles returns all immediate file entries.
// Returns a Files collection containing all regular files,
// excluding directories and other filesystem objects.
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
