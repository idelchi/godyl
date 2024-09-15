package remove

import (
	"fmt"

	"github.com/idelchi/godyl/internal/cache"
	"github.com/idelchi/godyl/pkg/path/folder"
)

func run(dir folder.Folder) error {
	cacheFile, err := cache.File(dir)
	if err != nil {
		return err
	}

	if err := cacheFile.Remove(); err != nil {
		return fmt.Errorf("removing cache: %w", err)
	}

	fmt.Printf("Cache file %q removed\n", cacheFile)

	return nil
}
