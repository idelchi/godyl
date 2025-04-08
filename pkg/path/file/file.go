package file

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/idelchi/godyl/pkg/utils"
	"gopkg.in/yaml.v3"
)

// File represents a file path as a string, providing methods for file operations.
type File string

// New creates a new File by joining the provided paths.
func New(paths ...string) File {
	return File(filepath.Clean(filepath.Join(paths...))) // .Normalized()
}

// Normalized converts the file path to use forward slashes.
func (f File) Normalized() File {
	return New(filepath.ToSlash(f.String()))
}

// Create creates a new file.
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

// Write writes the provided data to the file.
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

// Append appends the provided data to the file.
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

// Read reads the contents of the file and returns it as a byte slice.
func (f File) Read() ([]byte, error) {
	file, err := os.ReadFile(f.String())
	if err != nil {
		return nil, fmt.Errorf("reading file %q: %w", f, err)
	}

	return file, nil
}

// Remove deletes the file from the file system.
func (f File) Remove() error {
	if err := os.Remove(f.String()); err != nil {
		return fmt.Errorf("removing file %q: %w", f, err)
	}

	return nil
}

// Path returns the path of the File.
func (f File) Path() string {
	return f.String()
}

// WithoutExtension removes all recognized file extensions from the File.
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

// String returns the string representation of the File.
func (f File) String() string {
	return string(f)
}

// Base returns the base name of the File.
func (f File) Base() string {
	return filepath.Base(f.String())
}

// Join joins the File with the provided paths and returns a new File.
func (f File) Join(paths ...string) File {
	return New(append([]string{f.String()}, paths...)...)
}

// Expanded expands the file path in case of ~ and returns the expanded path.
func (f File) Expanded() File {
	return New(utils.ExpandHome(f.String()))
}

// Dir returns the directory of the file.
// If the file is a directory, it returns the path of the directory itself.
func (f File) Dir() string {
	if f.IsDir() {
		return f.String()
	}

	return filepath.Dir(f.String())
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

// Exists checks if the file exists in the file system.
func (f File) Exists() bool {
	_, err := os.Stat(f.String())

	return err == nil
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

// WriteYAML marshals the provided value to YAML and writes it to the file.
// It adds a newline at the end of the YAML content.
func (f File) WriteYAML(v any) error {
	// Marshal the struct to YAML
	data, err := yaml.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshaling to YAML: %w", err)
	}

	// Use the Write method to write the YAML data
	return f.Write(data)
}
