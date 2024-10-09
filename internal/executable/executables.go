package executable

import (
	"fmt"
	"path/filepath"
)

type Executables []Executable

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

func (es Executables) Find(dir string) (Executable, error) {
	for _, e := range es {
		executable, err := e.Find(dir)
		if err == nil {
			return executable, nil
		}
	}
	return Executable{}, fmt.Errorf("executable not found")
}

func (es Executables) Paths() (paths []string) {
	for _, e := range es {
		paths = append(paths, e.Path)
	}
	return paths
}

func (es Executables) SymlinksFor(exe Executable) error {
	return exe.Symlink(es.Paths())
}
