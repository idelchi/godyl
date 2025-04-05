package cache

import (
	"fmt"

	"github.com/idelchi/godyl/internal/cache/backend"
	"github.com/idelchi/godyl/internal/cache/cache"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/folder"
)

// File returns the cache file for the specified cache type.
func File(folder folder.Folder, cacheType string) file.File {
	switch cacheType {
	case "sqlite":
		return folder.WithFile("godyl.db")
	case "file":
		return folder.WithFile("godyl.json")
	default:
		return folder.WithFile("unknown")
	}
}

func New(folder folder.Folder, cacheType string) (*cache.Cache, error) {
	var backendType cache.Backend
	var err error
	switch cacheType {
	case "sqlite":
		backendType, err = backend.NewSQLite(File(folder, cacheType))
	case "file":
		backendType, err = backend.NewFile(File(folder, cacheType))
	default:
		return nil, fmt.Errorf("unsupported cache type: %s", cacheType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create %q backend: %w", cacheType, err)
	}

	// Create cache with file backend
	c := cache.New(backendType)

	return c, nil
}
