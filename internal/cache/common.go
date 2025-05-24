package cache

import (
	"fmt"
	"os"

	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/folder"
)

// File returns the cache file for the specified cache type.
// func File(folder folder.Folder) file.File {
// 	return folder.WithFile("godyl.json")
// }

// Existing returns the cache file from the specified folder.
func File(folder folder.Folder) (file.File, error) {
	cacheFile := folder.WithFile("godyl.json")

	if !cacheFile.Exists() {
		return cacheFile, fmt.Errorf("cache %w: %q", os.ErrNotExist, cacheFile)
	}

	return cacheFile, nil
}
