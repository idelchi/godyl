package path

import (
	"fmt"

	"github.com/idelchi/godyl/internal/cache"
	"github.com/idelchi/godyl/pkg/path/folder"
)

func run(dir folder.Folder) error {
	cacheFile, _ := cache.File(dir)

	fmt.Println(cacheFile)

	return nil
}
