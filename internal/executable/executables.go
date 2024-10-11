package executable

import (
	"fmt"
	"path/filepath"
)

// Executables is a list of Executable.
type Executables []Executable

// FromStrings creates a list of Executable from a list of strings.
func (Executables) FromStrings(dir string, executables ...string) Executables {
	e := Executables{}

	for _, executable := range executables {
		// Ignore empty strings
		if executable == "" {
			continue
		}

		e = append(e, Executable{Path: filepath.Join(dir, executable)})
	}

	return e
}

// Find searches for any of the executables in the given directory.
// The first executable found is returned.
func (es Executables) Find(dir string) (Executable, error) {
	for _, e := range es {
		executable, err := e.Find(dir)
		if err == nil {
			return executable, nil
		}
	}
	return Executable{}, fmt.Errorf("executables %v not found in %q", es.Paths(), dir)
}

// Paths returns a list of paths for all executables.
func (es Executables) Paths() (paths []string) {
	for _, e := range es {
		paths = append(paths, e.Path)
	}
	return paths
}

// SymlinksFor creates symlinks for all executables.
func (es Executables) SymlinksFor(exe Executable) error {
	return exe.Symlink(es.Paths())
}
