// Package data provides utilities for managing configuration and cache files and directories.
package data

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"

	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/folder"
)

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

// DownloadDir returns the temporary directory for downloads.
func DownloadDir() folder.Folder {
	return folder.New(os.TempDir())
}

// UserDataDir returns the user data directory for Godyl.
func UserDataDir() folder.Folder {
	dir, err := userDataDir()
	if err == nil {
		return folder.New(dir, "godyl")
	}

	if cd := ConfigDir(); cd.Path() != "." {
		return cd
	}

	return CacheDir()
}

// userDataDir returns the default root directory to use for user-specific
// data files. Users should create their own application-specific subdirectory
// within this one and use that.
//
// On Unix systems, it returns $XDG_DATA_HOME as specified by
// https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html if
// non-empty, else $HOME/.local/share.
// On Darwin, it returns $HOME/Library/Application Support.
// On Windows, it returns %AppData%.
// On Plan 9, it returns $home/lib.
//
// If the location cannot be determined (for example, $HOME is not defined) or
// the path in $XDG_DATA_HOME is relative, then it will return an error.
func userDataDir() (string, error) {
	var dir string

	switch runtime.GOOS {
	case "windows":
		dir = os.Getenv("AppData")
		if dir == "" {
			return "", errors.New("%AppData% is not defined")
		}

	case "darwin", "ios":
		dir = os.Getenv("HOME")
		if dir == "" {
			return "", errors.New("$HOME is not defined")
		}

		dir += "/Library/Application Support"

	case "plan9":
		dir = os.Getenv("home")
		if dir == "" {
			return "", errors.New("$home is not defined")
		}

		dir += "/lib"

	default: // Unix
		dir = os.Getenv("XDG_DATA_HOME")
		if dir == "" {
			dir = os.Getenv("HOME")
			if dir == "" {
				return "", errors.New("neither $XDG_DATA_HOME nor $HOME are defined")
			}

			dir += "/.local/share"
		} else if !filepath.IsAbs(dir) {
			return "", errors.New("path in $XDG_DATA_HOME is relative")
		}
	}

	return dir, nil
}

// GoDir returns the cache directory for Go installations.
// Defaults to temp directory if user cache directory cannot be determined.
func GoDir() folder.Folder {
	return CacheDir().Join("go")
}

// CacheFile returns the cache file from the specified folder.
func CacheFile(folder folder.Folder) file.File {
	return folder.WithFile("godyl.json")
}

// ConfigFile returns the first existing “godyl” configuration file it finds.
//
// Search order:
//  1. ./godyl.{yaml,yml}
//  2. $CONFIG_DIR/godyl.{yaml,yml}
//  3. ./godyl.yml
func ConfigFile() file.File {
	local := folder.New(".")

	for _, ext := range []string{"yaml", "yml"} {
		if f := local.WithFile("godyl").WithExtension(ext); f.Exists() {
			return f
		}
	}

	global := ConfigDir()

	for _, ext := range []string{"yaml", "yml"} {
		if f := global.WithFile("godyl").WithExtension(ext); f.Exists() {
			return f
		}
	}

	return global.WithFile("godyl.yml")
}

// CreateUniqueDirIn creates a unique directory in the specified path.
// If no path is specified, it uses the system temporary directory.
func CreateUniqueDirIn(paths ...string) (folder.Folder, error) {
	if len(paths) == 0 {
		paths = []string{DownloadDir().Path()}
	}

	return folder.CreateRandomInDir(folder.New(paths...).Path(), Prefix())
}

// Prefix returns the prefix used for Godyl temporary directories.
func Prefix() string {
	return "godyl-*"
}
