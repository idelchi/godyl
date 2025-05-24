package path

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

	fmt.Println(cacheFile)

	return nil
}
