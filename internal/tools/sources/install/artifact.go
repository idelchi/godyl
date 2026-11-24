package install

import (
	"errors"
	"fmt"
	"io/fs"

	"github.com/idelchi/godyl/internal/data"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/folder"
)

// Artifact owns a staged download until Remove is called.
// Find may share an artifact; extraction needs exclusive access to adjust staged permissions.
type Artifact struct {
	// Destination is the loose file or extracted directory returned by go-getter.
	Destination file.File
	// dir is the private parent, whose permissions cannot be changed by archive headers.
	dir folder.Folder
}

// Acquire downloads an artifact into private staging. The caller owns its cleanup.
func Acquire(d Data) (Artifact, error) {
	dir, err := data.CreateUniqueDirIn()
	if err != nil {
		return Artifact{}, fmt.Errorf("creating download directory: %w", err)
	}

	artifact := Artifact{dir: dir}

	destination, err := downloadInto(d, artifact.contents())
	if err != nil {
		return Artifact{}, errors.Join(err, artifact.Remove())
	}

	artifact.Destination = destination

	return artifact, nil
}

// Remove releases the artifact, including read-only directories created by tar headers.
func (a Artifact) Remove() error {
	if err := a.dir.Remove(); !errors.Is(err, fs.ErrPermission) {
		return err
	}

	const ownerOnly = 0o700

	err := a.dir.Walk(func(path file.File, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			return path.Chmod(ownerOnly)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return a.dir.Remove()
}

// contents is the directory into which go-getter writes.
func (a Artifact) contents() folder.Folder {
	return a.dir.Join("artifact")
}
