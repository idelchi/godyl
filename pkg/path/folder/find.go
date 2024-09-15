package folder

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/idelchi/godyl/pkg/path/file"
)

// FindFile searches for a file matching all given criteria.
// Recursively searches the directory tree and returns the first matching file.
// Returns ErrNotFound if no file matches all criteria.
func (f Folder) FindFile(criteria ...CriteriaFunc) (file.File, error) {
	root := f.AsFile()

	var foundPath file.File

	err := filepath.WalkDir(root.String(), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walking folder %q: %w", root, err)
		}

		if d.IsDir() {
			return nil // skip directories
		}

		current := file.New(path)

		relPath, err := current.RelativeTo(root)
		if err != nil {
			return fmt.Errorf("getting relative path: %w", err)
		}

		for _, criterion := range criteria {
			matches, err := criterion(relPath)
			if err != nil {
				return fmt.Errorf("evaluating criterion: %w", err)
			}

			if !matches {
				return nil
			}
		}

		foundPath = root.Join(relPath.String())

		return filepath.SkipAll
	})
	if err != nil {
		return file.File(""), fmt.Errorf("walking folder %q: %w", root, err)
	}

	if foundPath == "" {
		return file.File(""), fmt.Errorf("%w: no file found matching all criteria in folder %q", ErrNotFound, root)
	}

	return foundPath, nil
}
