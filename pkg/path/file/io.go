package file

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
)

// createFolder creates the directory for this file if it does not exist.
func (f File) createFolder() error {
	const perm = 0o755

	if err := os.MkdirAll(f.Dir(), perm); err != nil {
		return fmt.Errorf("creating directory %q: %w", f, err)
	}

	return nil
}

// CreateRandomInDir creates a uniquely named file.
// Creates a file with a random name inside the specified directory.
// Use empty string for directory to create in system temp directory.
// Pattern is used as a prefix for the random file name.
func CreateRandomInDir(dir, pattern string) (File, error) {
	file, err := os.CreateTemp(dir, pattern)
	if err != nil {
		return File(""), fmt.Errorf("creating temporary file in %q: %w", dir, err)
	}

	if err := file.Close(); err != nil {
		return File(""), fmt.Errorf("closing temporary file %q: %w", file.Name(), err)
	}

	return New(file.Name()), nil
}

// Create creates a new empty file at this path.
// Returns an error if the file cannot be created.
func (f File) Create() error {
	if err := f.createFolder(); err != nil {
		return err
	}

	file, err := os.Create(f.Path())
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

	file, err := os.OpenFile(f.Path(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return nil, fmt.Errorf("opening file %q for writing: %w", f, err)
	}

	return file, nil
}

// Write stores binary data in the file.
// Creates or truncates the file before writing.
func (f File) Write(data []byte) (err error) {
	file, err := f.OpenForWriting()
	if err != nil {
		return err
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			err = errors.Join(err, fmt.Errorf("closing file %q: %w", f, cerr))
		}
	}()

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("writing to file %q: %w", f, err)
	}

	return nil
}

// Open opens the file for reading and returns a pointer to the os.File object.
// The user must close the file after use.
func (f File) Open() (*os.File, error) {
	file, err := os.Open(f.Path())
	if err != nil {
		return nil, fmt.Errorf("opening file %q: %w", f, err)
	}

	return file, nil
}

// Read retrieves the entire contents of the file.
// Returns the file contents as a byte slice.
func (f File) Read() ([]byte, error) {
	file, err := os.ReadFile(f.Path())
	if err != nil {
		return nil, fmt.Errorf("reading file %q: %w", f, err)
	}

	return file, nil
}

// ReadString retrieves the entire contents of the file as a string.
func (f File) ReadString() (string, error) {
	data, err := f.Read()
	if err != nil {
		return "", fmt.Errorf("reading file %q as string: %w", f, err)
	}

	return string(data), nil
}

// Lines retrieves the file contents as a slice of strings.
func (f File) Lines() ([]string, error) {
	data, err := f.Read()
	if err != nil {
		return nil, fmt.Errorf("reading file %q: %w", f, err)
	}

	// Count lines for pre-allocation
	lineCount := strings.Count(string(data), "\n") + 1

	returns := make([]string, 0, lineCount)

	for line := range strings.SplitSeq(string(data), "\n") {
		returns = append(returns, strings.TrimRight(line, "\r"))
	}

	return returns, nil
}

// Remove deletes the file from the filesystem.
// Returns an error if the file cannot be deleted.
// Silently ignores the error if the file does not exist.
func (f File) Remove() error {
	if err := os.Remove(f.Path()); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("removing file %q: %w", f, err)
	}

	return nil
}

// Chmod modifies the file's permission bits.
// Takes a fs.FileMode parameter specifying the new permissions.
func (f *File) Chmod(mode fs.FileMode) error {
	if err := os.Chmod(f.Path(), mode); err != nil {
		return fmt.Errorf("changing permissions of file %q: %w", f, err)
	}

	return nil
}

// Copy duplicates the file to a new location.
// It copies contents and preserves the original file's permissions.
// Timestamps and ownership are not preserved for cross-platform compatibility.
func (f File) Copy(dest File) error {
	if dest.Absolute().Path() == f.Absolute().Path() {
		return fmt.Errorf("source (%q) and destination (%q) are identical", f, dest)
	}

	// Delete existing destination file if it exists
	if err := dest.Remove(); err != nil {
		return fmt.Errorf("removing existing destination file %q: %w", dest, err)
	}

	// Open source file
	src, err := f.Open()
	if err != nil {
		return fmt.Errorf("opening source file: %w", err)
	}
	defer src.Close()

	// Stat source to get mode
	srcInfo, err := os.Stat(f.Path())
	if err != nil {
		return fmt.Errorf("statting source file: %w", err)
	}

	// Open destination with same mode
	dst, err := os.OpenFile(dest.Path(), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode().Perm())
	if err != nil {
		return fmt.Errorf("opening destination file: %w", err)
	}
	defer dst.Close()

	// Copy contents
	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("copying file contents: %w", err)
	}

	return nil
}

// Copies duplicates the file to a new location.
// It copies contents and preserves the original file's permissions.
// Timestamps and ownership are not preserved for cross-platform compatibility.
func (f File) Copies(copies ...File) error {
	for _, copy := range copies {
		if copy.Absolute().Path() == f.Absolute().Path() {
			continue
		}

		if err := f.Copy(copy); err != nil {
			return fmt.Errorf("copying file %q to %q: %w", f, copy, err)
		}
	}

	return nil
}

// Hardlinks creates a hard link at the destination path pointing to this file.
// Returns an error if the operation is not permitted, the files are on different devices,
// or the source file is not a regular file.
func (f File) Hardlinks(hardlinks ...File) error {
	for _, hardlink := range hardlinks {
		if hardlink.Absolute().Path() == f.Absolute().Path() {
			continue
		}

		if err := hardlink.Remove(); err != nil {
			return fmt.Errorf("removing existing hardlink %q: %w", hardlink, err)
		}

		err := os.Link(f.Absolute().Path(), hardlink.Absolute().Path())
		if err != nil && !os.IsExist(err) {
			return fmt.Errorf("creating hardlink for %q: %w", hardlink, err)
		}
	}

	return nil
}

// Softlinks creates symbolic links to this file on Unix-like systems.
// Takes multiple target paths and creates a symlink at each location.
// Skips existing symlinks without error, but returns an error if
// symlink creation fails for any other reason. Not available on Windows.
func (f File) Softlinks(softlinks ...File) error {
	for _, softlink := range softlinks {
		if softlink.Absolute().Path() == f.Absolute().Path() {
			continue
		}

		if err := softlink.Remove(); err != nil {
			return fmt.Errorf("removing existing softlink %q: %w", softlink, err)
		}

		err := os.Symlink(f.Absolute().Path(), softlink.Absolute().Path())
		if err != nil && !os.IsExist(err) {
			return fmt.Errorf("creating softlink for %q: %w", softlink, err)
		}
	}

	return nil
}

// Links tries to create symbolic links to this file at the given targets.
// If symlinking fails (e.g., unsupported platform, permissions, or FS restrictions),
// it attempts a hard link. If that also fails, it falls back to copying the file.
//
// This function is a convenience helper for creating portable references
// to a file without requiring the caller to explicitly handle platform quirks.
func (f File) Links(links ...File) error {
	for _, link := range links {
		if link.Absolute().Path() == f.Absolute().Path() {
			continue
		}

		// Try symlink
		if err := f.Softlinks(link); err == nil {
			continue
		}
		// Try hard link
		if err := f.Hardlinks(link); err == nil {
			continue
		}

		// Fallback to full copy
		if err := f.Copies(link); err != nil {
			return fmt.Errorf("link fallback copy for %q: %w", link, err)
		}
	}

	return nil
}
