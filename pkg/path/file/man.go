package file

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"

	"github.com/idelchi/godyl/pkg/utils"
)

// Matches checks if the file path matches the given (extended glob) pattern.
func (f File) Matches(pattern string) (bool, error) {
	return doublestar.Match(pattern, f.String())
}

// RelativeTo computes the relative path from base to this file.
// Returns an error if the relative path cannot be computed.
func (f File) RelativeTo(base File) (File, error) {
	rel, err := filepath.Rel(base.String(), f.String())
	if err != nil {
		return f, fmt.Errorf("getting relative path from %q to %q: %w", base, f, err)
	}

	return New(rel), nil
}

// Path returns the string representation of the file path.
func (f File) Path() string {
	return f.String()
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

// Absolute returns the absolute path of the file.
func (f File) Absolute() File {
	absPath, err := filepath.Abs(f.String())
	if err != nil {
		return f
	}

	return New(absPath)
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

// WithExtension returns a new File with the specified suffix added to the original file name.
func (f File) WithExtension(extension string) File {
	return New(fmt.Sprintf("%s.%s", f.WithoutExtension().Path(), extension))
}

// WithoutExtension returns the path without file extensions.
func (f File) WithoutExtension() File {
	return New(strings.TrimSuffix(f.String(), "."+f.Extension()))
}

// HasExtension checks if the file has an extension.
func (f File) HasExtension() bool {
	return f.Extension() != ""
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
