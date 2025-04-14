package file

import (
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/idelchi/godyl/pkg/utils"
	"gopkg.in/yaml.v3"
)

// File represents a filesystem path with associated operations.
// Provides methods for file manipulation, path operations, and
// common file system tasks.
type File string

// New creates a File from one or more path components.
// Joins the paths using filepath.Join and normalizes the result.
func New(paths ...string) File {
	return File(filepath.Clean(filepath.Join(paths...))).Normalized()
}

// Normalized returns the path with forward slashes.
// Converts backslashes to forward slashes for consistency.
func (f File) Normalized() File {
	return File(filepath.ToSlash(f.String()))
}

// RelativeTo computes the relative path from this file to a base path.
// Returns an error if the relative path cannot be computed.
func (f File) RelativeTo(base File) (File, error) {
	file, err := filepath.Rel(f.String(), base.String())
	if err != nil {
		return f, fmt.Errorf("getting relative path: %w", err)
	}

	return New(file), nil
}

// Create creates a new empty file at this path.
// Returns an error if the file cannot be created.
func (f File) Create() error {
	file, err := os.Create(f.String())
	if err != nil {
		return fmt.Errorf("creating file %q: %w", f, err)
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("closing file %q: %w", f, err)
	}

	return nil
}

// OpenForWriting opens the file for writing and returns a pointer to the os.File object.
// If the file doesn't exist, it will be created.
// If it exists, it will be truncated.
// The user must close the file after use.
func (f File) OpenForWriting() (*os.File, error) {
	const perm = 0o600

	file, err := os.OpenFile(f.String(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return nil, fmt.Errorf("opening file %q for writing: %w", f, err)
	}

	return file, nil
}

// Write stores binary data in the file.
// Creates or truncates the file before writing.
func (f File) Write(data []byte) error {
	file, err := f.OpenForWriting()
	if err != nil {
		return fmt.Errorf("opening file %q for writing: %w", f, err)
	}

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("writing to file %q: %w", f, err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("closing file %q: %w", f, err)
	}
	return nil
}

// OpenForAppend opens the file for appending and returns a pointer to the os.File object.
// If the file doesn't exist, it will be created.
// The user must close the file after use.
func (f File) OpenForAppend() (*os.File, error) {
	const perm = 0o600

	file, err := os.OpenFile(f.String(), os.O_WRONLY|os.O_CREATE|os.O_APPEND, perm)
	if err != nil {
		return nil, fmt.Errorf("opening file %q for appending: %w", f, err)
	}

	return file, nil
}

// Append adds binary data to the end of the file.
// Creates the file if it doesn't exist.
func (f File) Append(data []byte) error {
	file, err := f.OpenForAppend()
	if err != nil {
		return err
	}

	if _, err := file.Write(data); err != nil {
		file.Close()
		return fmt.Errorf("appending to file %q: %w", f, err)
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("closing file %q: %w", f, err)
	}

	return nil
}

// Open opens the file for reading and returns a pointer to the os.File object, or an error.
// The user must close the file after use.
func (f File) Open() (*os.File, error) {
	file, err := os.Open(f.String())
	if err != nil {
		return nil, fmt.Errorf("opening file %q: %w", f, err)
	}

	return file, nil
}

// Read retrieves the entire contents of the file.
// Returns the file contents as a byte slice.
func (f File) Read() ([]byte, error) {
	file, err := os.ReadFile(f.String())
	if err != nil {
		return nil, fmt.Errorf("reading file %q: %w", f, err)
	}

	return file, nil
}

// Remove deletes the file from the filesystem.
// Returns an error if the file cannot be deleted.
func (f File) Remove() error {
	if err := os.Remove(f.String()); err != nil {
		return fmt.Errorf("removing file %q: %w", f, err)
	}

	return nil
}

// Path returns the string representation of the file path.
func (f File) Path() string {
	return f.String()
}

// WithoutExtension returns the path without file extensions.
// Handles compound extensions (e.g., .tar.gz) and unknown extensions.
func (f File) WithoutExtension() File {
	result := f
	for ext := result.Extension(); ext != None && ext != Other; ext = result.Extension() {
		result = File(strings.TrimSuffix(strings.ToLower(result.String()), "."+strings.ToLower(ext.String())))
	}

	if ext := result.Extension(); ext == Other {
		// Trim everything after and including the last dot
		result = File(strings.TrimSuffix(result.String(), filepath.Ext(result.String())))
	}

	return result
}

// String returns the file path as a string.
func (f File) String() string {
	return string(f)
}

// Base returns the last component of the file path.
func (f File) Base() string {
	return filepath.Base(f.String())
}

// Join combines this path with additional components.
// Returns a new File with the combined path.
func (f File) Join(paths ...string) File {
	return New(append([]string{f.String()}, paths...)...)
}

// Expanded resolves home directory references.
// Replaces ~ with the user's home directory path.
func (f File) Expanded() File {
	return New(utils.ExpandHome(f.String()))
}

// Dir returns the containing directory path.
// Returns the path itself if it's a directory,
// otherwise returns the parent directory path.
func (f File) Dir() string {
	if f.IsDir() {
		return f.String()
	}

	return filepath.Dir(f.String())
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

// Chmod modifies the file's permission bits.
// Takes a fs.FileMode parameter specifying the new permissions.
func (f *File) Chmod(mode fs.FileMode) error {
	return os.Chmod(f.String(), mode)
}

// Copy duplicates the file to a new location.
// Creates the destination file with execute permissions (0755)
// and copies all content from the source file.
func (f File) Copy(other File) error {
	// Open the source file
	source, err := f.Open()
	if err != nil {
		return fmt.Errorf("opening source file: %w", err)
	}
	defer source.Close()

	// Create the destination file
	destination, err := other.OpenForWriting()
	if err != nil {
		return fmt.Errorf("opening destination file: %w", err)
	}
	defer destination.Close()

	// Copy the contents of the source file to the destination file
	_, err = io.Copy(destination, source)
	if err != nil {
		return fmt.Errorf("copying file: %w", err)
	}

	const perm = 0o755

	// Set permissions on the destination file (executable permission)
	if err := destination.Chmod(perm); err != nil {
		return fmt.Errorf("setting permissions: %w", err)
	}

	return nil
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

// Extension returns the file's extension type.
// Maps the extension to a predefined constant (e.g., EXE, TAR, ZIP).
func (f File) Extension() Extension {
	ext := filepath.Ext(f.String())

	switch strings.ToLower(ext) {
	case ".exe":
		return EXE
	case ".gz":
		return GZ
	case ".zip":
		return ZIP
	case ".tar":
		return TAR
	case "":
		return None
	default:
		return Other
	}
}

// WriteYAML serializes data to YAML format and writes it.
// Marshals the provided value to YAML and writes to the file.
func (f File) WriteYAML(v any) error {
	// Marshal the struct to YAML
	data, err := yaml.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshaling to YAML: %w", err)
	}

	// Use the Write method to write the YAML data
	return f.Write(data)
}

// Unescape decodes URL-escaped characters in the path.
// Returns the original path if unescaping fails.
func (f File) Unescape() File {
	// Replace escaped characters with their original representation
	unescapedPath, err := url.QueryUnescape(f.String())
	if err != nil {
		return f
	}

	return File(unescapedPath)
}
