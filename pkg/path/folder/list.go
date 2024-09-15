package folder

import (
	"fmt"
	"os"
	"slices"

	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/files"
)

// ListFolders returns all immediate subdirectories of the folder.
// It excludes files and other non-directory entries.
func (f Folder) ListFolders() (folders []Folder, err error) {
	entries, err := os.ReadDir(f.String())
	if err != nil {
		return nil, fmt.Errorf("reading directory %q: %w", f, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			subFolder := New(f.String(), entry.Name())
			folders = append(folders, subFolder)
		}
	}

	slices.Sort(folders)

	return folders, nil
}

// ListFiles returns all immediate regular files in the folder.
// It excludes directories and other non-file entries.
func (f Folder) ListFiles() (files files.Files, err error) {
	entries, err := os.ReadDir(f.String())
	if err != nil {
		return nil, fmt.Errorf("reading directory %q: %w", f, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			file := file.New(f.String(), entry.Name())
			files = append(files, file)
		}
	}

	slices.Sort(files)

	return files, nil
}
