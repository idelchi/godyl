package cache

import (
	"fmt"

	"github.com/idelchi/godyl/internal/cache/backend"
	"github.com/idelchi/godyl/internal/cache/cache"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/folder"
)

// File returns the cache file for the specified cache type.
func File(folder folder.Folder) file.File {
	return folder.WithFile("godyl.json")
}

func New(folder folder.Folder) (*cache.Cache, error) {
	backendType, err := backend.NewFile(File(folder))
	if err != nil {
		return nil, fmt.Errorf("failed to create cache backend: %w", err)
	}

	return cache.New(backendType), nil
}
