// Package tmp provides utilities for temporary file and directory management.
package tmp

import (
	"os"

	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/folder"
)

// ConfigFile returns the first existing “godyl” configuration file it finds.
//
// Search order:
//  1. $CONFIG_DIR/godyl.yaml
//  2. $CONFIG_DIR/godyl.yml
//  3. ./godyl.yml (project root)
//
// If none exist, it falls back to $CONFIG_DIR/godyl.yml, allowing callers
// to create or write it later.
func ConfigFile() file.File {
	base := ConfigDir().WithFile("godyl")

	// Check config directory for YAML/YML variants.
	for _, ext := range []string{"yaml", "yml"} {
		if f := base.WithExtension(ext); f.Exists() {
			return f
		}
	}

	// Fallback: project-local godyl.yml
	if local := file.New("godyl.yml"); local.Exists() {
		return local
	}

	// Default path when nothing exists yet.
	return base.WithExtension("yml")
}

// ConfigDir returns the config directory for Godyl.
func ConfigDir() folder.Folder {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return folder.New(".")
	}

	return folder.New(configDir, "godyl")
}

// CacheDir returns the cache directory for Godyl.
func CacheDir() folder.Folder {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return folder.New(os.TempDir(), "godyl")
	}

	return folder.New(cacheDir, "godyl")
}

// CacheFile returns the cache file from the specified folder.
func CacheFile(folder folder.Folder) file.File {
	return folder.WithFile("godyl.json")
}

// DownloadDir returns the temporary directory for Godyl.
// Optionally pass in subdirectories to create a path within the Godyl directory.
func DownloadDir(paths ...string) folder.Folder {
	return folder.New(os.TempDir()).Join(paths...)
}

// GoDir returns the temporary directory for Go installations.
func GoDir() folder.Folder {
	return CacheDir().Join("go")
}

// CreateRandomDir creates a random directory in the Godyl temporary directory.
func CreateRandomDir() (folder.Folder, error) {
	// Create a random temporary directory for Godyl
	return folder.CreateRandomInDir(DownloadDir().Path(), Prefix())
}

// CreateRandomDirIn creates a random directory in the specified directory.
func CreateRandomDirIn(dir folder.Folder) (folder.Folder, error) {
	// Create a random temporary directory for Godyl
	return folder.CreateRandomInDir(dir.Path(), Prefix())
}

// Prefix returns the prefix used for Godyl temporary directories.
func Prefix() string {
	return "godyl-*"
}
