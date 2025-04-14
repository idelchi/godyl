package file

import (
	"fmt"
)

// Symlink emulates symbolic links on Windows by file copying.
// Takes multiple target paths and creates a copy at each location.
// Skips copying if source and target paths are identical.
// Returns an error if any copy operation fails. This is a Windows-
// specific implementation that works around symlink limitations.
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
