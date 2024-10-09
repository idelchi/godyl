package folder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Folder string

func (f *Folder) IsSet() bool {
	return *f != ""
}

func (f *Folder) IsParentOf(child Folder) bool {
	return strings.HasPrefix(child.Path(), f.Path())
}

func (f *Folder) Expand() error {
	if strings.HasPrefix(f.Path(), "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("getting user home directory: %w", err)
		}
		*f = Folder(filepath.Join(homeDir, f.Path()[2:]))
	}

	return nil
}

func (f Folder) Path() string {
	return string(f)
}

func (f Folder) Exists() bool {
	_, err := os.Stat(f.Path())
	return err == nil
}

func (f Folder) Create() error {
	return os.MkdirAll(f.Path(), 0o755)
}

func (f Folder) Name() string {
	return filepath.Base(f.Path())
}

func (f *Folder) CreateRandomInTempDir() error {
	name, err := os.MkdirTemp("", "godyl-*")
	*f = Folder(name)

	return err
}

func (f *Folder) CreateInTempDir() error {
	name := filepath.Join(os.TempDir(), f.Name())

	err := os.Mkdir(name, 0o755)

	*f = Folder(name)

	return err
}

func (f Folder) Remove() error {
	return os.RemoveAll(f.Path())
}
