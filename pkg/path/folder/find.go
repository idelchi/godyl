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
	files, err := f.FindFiles(criteria...)

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
func (f Folder) FindFiles(criteria ...CriteriaFunc) (files.Files, error) {
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

		return nil
	})
	if err != nil {
		return files.Files{}, fmt.Errorf("walking folder %q: %w", root, err)
	}

	slices.Sort(found)

	return found, nil
}

// Glob returns all files matching the given pattern within the folder.
// Uses filepath.Glob semantics (non-recursive, standard glob syntax).
func (f Folder) Glob(pattern string) (files.Files, error) {
	matches, err := filepath.Glob(filepath.Join(f.Path(), pattern))
	if err != nil {
		return nil, fmt.Errorf("glob %q in %q: %w", pattern, f, err)
	}

	result := make(files.Files, 0, len(matches))
	for _, m := range matches {
		result = append(result, file.New(m))
	}

	return result, nil
}

// Walk traverses the directory tree rooted at f, calling fn for each entry.
// The error parameter from filepath.WalkDir is passed through — return nil to
// skip the error, or return it to abort the walk. Return filepath.SkipDir to
// skip a directory.
func (f Folder) Walk(fn func(path file.File, d fs.DirEntry, err error) error) error {
	return filepath.WalkDir(f.Path(), func(p string, d fs.DirEntry, err error) error {
		return fn(file.New(p), d, err)
	})
}
