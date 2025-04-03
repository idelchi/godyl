package tmp

import (
	"os"
	"path/filepath"

	"github.com/idelchi/godyl/pkg/folder"
)

// GodylDir returns the temporary directory for Godyl.
// Optionally pass in subdirectories to create a path within the Godyl directory.
func GodylDir(paths ...string) folder.Folder {
	return folder.New(os.TempDir(), ".godyl", filepath.Join(paths...))
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
