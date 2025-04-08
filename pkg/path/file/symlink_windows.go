package file

import (
	"fmt"
)

// Symlink creates symbolic links by copying the content of the File to each of the provided symlink Files.
// It skips the operation if the symlink has the same name as the original File.
// Returns an error if any of the copy operations fail.
func (f File) Symlink(symlinks ...File) error {
	for _, symlink := range symlinks {
		if symlink.Path() == f.Path() {
			continue
		}

		if err := f.Copy(symlink); err != nil {
			return fmt.Errorf("copying %q to %q: %w", f, symlink, err)
		}
	}

	return nil
}
