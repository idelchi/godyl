package file

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"

	"github.com/idelchi/godyl/pkg/generic"
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

// ExpandedX resolves home directory references.
// Replaces ~ with the user's home directory path.
func (f File) ExpandedX() File {
	return New(generic.ExpandHome(f.String()))
}

// Expanded resolves environment variables including `~` home directory references.
func (f File) Expanded() File {
	return New(os.ExpandEnv(generic.ExpandHome(f.String())))
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

// MakeExecutable sets the file as executable for user, group, and others.
func (f File) MakeExecutable() error {
	// Get the current file info
	info, err := f.Info()
	if err != nil {
		return fmt.Errorf("getting file info: %w", err)
	}

	const executableMask = 0o111 // User, Group, and Others execute permissions

	// Set the executable bit for user, group, and others
	newMode := info.Mode() | executableMask

	// Change the file mode
	if err := f.Chmod(newMode); err != nil {
		return fmt.Errorf("changing file mode: %w", err)
	}

	return nil
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

// WithoutFolder strips the leading folder path from f if it matches as a full
// path prefix (on a segment boundary). It never leaves a leading "/" behind.
func (f File) WithoutFolder(prefix string) File {
	fp := strings.TrimPrefix(f.Path(), "./")

	// Normalize the prefix: slashes, drop leading "./", and strip trailing "/"
	p := filepath.ToSlash(prefix)

	p = strings.TrimPrefix(p, "./")
	p = strings.Trim(p, "/") // handles both "" and trailing "/"

	if p == "" {
		return f
	}

	// Only strip when the prefix is a full segment: "p/"
	if rest, ok := strings.CutPrefix(fp, p+"/"); ok {
		return New(rest)
	}

	// Exact dir match or no match -> leave unchanged
	return f
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
