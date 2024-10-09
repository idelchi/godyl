package executable

import (
	"fmt"
)

func (e Executable) Symlink(symlinks []string) error {
	for _, symlink := range symlinks {
		if symlink == e.Path {
			continue
		}

		if err := e.Copy(symlink); err != nil {
			return fmt.Errorf("copying %q to %q: %w", e.Path, symlink, err)
		}
	}

	return nil
}
