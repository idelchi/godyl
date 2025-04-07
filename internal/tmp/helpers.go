package tmp

import (
	"os"

	"github.com/idelchi/godyl/pkg/path/folder"
)

// CacheDir returns the cache directory for Godyl.
func CacheDir() folder.Folder {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return folder.New(os.TempDir(), "godyl")
	}

	return folder.New(cacheDir, "godyl")
}

// DownloadDir returns the download directory for Godyl.
func DownloadDir() folder.Folder {
	return folder.New(os.TempDir(), "godyl")
}

// GodylDir returns the temporary directory for Godyl.
// Optionally pass in subdirectories to create a path within the Godyl directory.
func GodylDir(paths ...string) folder.Folder {
	return DownloadDir().Join(paths...)
}

// GodylCreateRandomDir creates a random directory in the Godyl temporary directory.
func GodylCreateRandomDir() (folder.Folder, error) {
	// Create a random temporary directory for Godyl
	return folder.CreateRandomInDir(GodylDir().Path(), Prefix())
}

// GodylCreateRandomDirIn creates a random directory in the specified directory.
func GodylCreateRandomDirIn(dir folder.Folder) (folder.Folder, error) {
	// Create a random temporary directory for Godyl
	return folder.CreateRandomInDir(dir.Path(), Prefix())
}

// Prefix returns the prefix used for Godyl temporary directories.
func Prefix() string {
	return "godyl-*"
}
