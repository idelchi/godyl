package folder

import (
	"errors"

	"github.com/idelchi/godyl/pkg/path/file"
)

// Folder represents a filesystem directory path.
// Provides methods for directory operations including creation removal, path manipulation, and file searching.
type Folder string

// CriteriaFunc defines a file matching predicate.
// Returns true if a file matches the criteria, false otherwise.
type CriteriaFunc func(file.File) (bool, error)

// ErrNotFound indicates a file matching the search criteria was not found.
var ErrNotFound = errors.New("file not found")
