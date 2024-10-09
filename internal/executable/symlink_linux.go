package executable

import (
	"fmt"
	"os"
)

func (e Executable) Symlink(symlinks []string) error {
	for _, symlink := range symlinks {
		if symlink == e.Path {
			continue
		}

		err := os.Symlink(e.Path, symlink)
		if err != nil && !os.IsExist(err) {
			return fmt.Errorf("creating symlink for %q: %w", symlink, err)
		}
	}

	return nil
}
