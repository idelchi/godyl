package folder

import (
	"cmp"
	"fmt"
	"io/fs"
	"path/filepath"
	"slices"

	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/files"
)

// FindFile searches for a file matching all given criteria.
// Recursively searches the directory tree and returns the first matching file, with the shortest path.
// Returns ErrNotFound if no file matches all criteria.
func (f Folder) FindFile(criteria ...CriteriaFunc) (file.File, error) {
	files, err := f.FindFiles(true, criteria...)

	if len(files) > 0 {
		return slices.MinFunc(files, func(a, b file.File) int {
			return cmp.Compare(len(a), len(b))
		}), err
	}

	return file.File(""), fmt.Errorf("%w: no file found matching all criteria in folder %q", ErrNotFound, f.AsFile())
}

// FindFiles searches for files matching all given criteria.
// Recursively searches the directory tree and returns all matching files.
// Returns an empty slice if no files match all criteria.
func (f Folder) FindFiles(firstOnly bool, criteria ...CriteriaFunc) (files.Files, error) {
	root := f.AsFile()

	var found files.Files

	err := filepath.WalkDir(root.Path(), func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("walking folder %q: %w", root, walkErr)
		}

		if d.IsDir() || !d.Type().IsRegular() {
			return nil
		}

		current := file.New(path)

		relPath, err := current.RelativeTo(root.Path())
		if err != nil {
			return fmt.Errorf("getting relative path: %w", err)
		}

		for _, criterion := range criteria {
			match, err := criterion(relPath)
			if err != nil {
				return fmt.Errorf("evaluating criterion: %w", err)
			}

			if !match {
				return nil
			}
		}

		found = append(found, root.Join(relPath.Path()))

		if firstOnly {
			return filepath.SkipAll
		}

		return nil
	})
	if err != nil {
		return files.Files{}, fmt.Errorf("walking folder %q: %w", root, err)
	}

	slices.Sort(found)

	return found, nil
}
