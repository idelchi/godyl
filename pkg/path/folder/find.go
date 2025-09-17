package folder

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"slices"

	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/files"
)

// FindFile searches for a file matching all given criteria.
// Recursively searches the directory tree and returns the first matching file.
// Returns ErrNotFound if no file matches all criteria.
func (f Folder) FindFile(criteria ...CriteriaFunc) (file.File, error) {
	files, err := f.FindFiles(true, criteria...)

	if len(files) > 0 {
		return files[0], err
	}

	return file.File(""), fmt.Errorf("%w: no file found matching all criteria in folder %q", ErrNotFound, f.AsFile())

	// root := f.AsFile()

	// var foundPath file.File

	// err := filepath.WalkDir(root.String(), func(path string, d fs.DirEntry, err error) error {
	// 	if err != nil {
	// 		return fmt.Errorf("walking folder %q: %w", root, err)
	// 	}

	// 	if d.IsDir() {
	// 		return nil
	// 	}

	// 	current := file.New(path)

	// 	relPath, err := current.RelativeTo(root)
	// 	if err != nil {
	// 		return fmt.Errorf("getting relative path: %w", err)
	// 	}

	// 	for _, criterion := range criteria {
	// 		matches, err := criterion(relPath)
	// 		if err != nil {
	// 			return fmt.Errorf("evaluating criterion: %w", err)
	// 		}

	// 		if !matches {
	// 			return nil
	// 		}
	// 	}

	// 	foundPath = root.Join(relPath.String())

	// 	return filepath.SkipAll
	// })
	// if err != nil {
	// 	return file.File(""), fmt.Errorf("walking folder %q: %w", root, err)
	// }

	// if foundPath == "" {
	// 	return file.File(""), fmt.Errorf("%w: no file found matching all criteria in folder %q", ErrNotFound, root)
	// }

	// return foundPath, nil
}

// FindFiles searches for files matching all given criteria.
// Recursively searches the directory tree and returns all matching files.
// Returns an empty slice if no files match all criteria.
func (f Folder) FindFiles(findFirst bool, criteria ...CriteriaFunc) (files.Files, error) {
	root := f.AsFile()

	var found files.Files

	err := filepath.WalkDir(root.String(), func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("walking folder %q: %w", root, walkErr)
		}

		if d.IsDir() || !d.Type().IsRegular() {
			return nil
		}

		current := file.New(path)

		relPath, err := current.RelativeTo(root)
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

		found = append(found, root.Join(relPath.String()))

		if findFirst {
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
