package install

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/folder"
)

// extract merges an archive or installs a loose file under its getter-selected name.
func (a Artifact) extract(output folder.Folder) error {
	if a.Destination.IsDir() {
		return extractDir(folder.FromFile(a.Destination), output)
	}

	name, err := a.Destination.RelativeTo(a.contents().Path())
	if err != nil {
		return err
	}

	target := output.WithFile(name.Path())
	if err := folder.FromFile(target).Create(); err != nil {
		return err
	}

	target, info, err := extractTarget(target)
	if err != nil {
		return err
	}

	// Unlike archive members, a raw HTTP download retains an existing file's permissions.
	if info == nil {
		return a.Destination.CopyAtomic(target, nil)
	}

	return a.Destination.CopyAtomic(target, func(tmp file.File) error {
		return tmp.Chmod(info.Mode().Perm())
	})
}

// extractDir merges staged entries without changing permissions on existing output directories.
func extractDir(source, output folder.Folder) error {
	created := !output.Exists()

	info, err := source.Info()
	if err != nil {
		return err
	}

	// An archive can remove the owner's read/search bits. Only this installer uses this staging tree.
	const readSearch = 0o500

	if info.Mode().Perm()&readSearch != readSearch {
		if err := source.Chmod(info.Mode().Perm() | readSearch); err != nil {
			return err
		}
	}

	if err := output.Create(); err != nil {
		return err
	}

	entries, err := source.List()
	if err != nil {
		return err
	}

	for _, entry := range entries {
		name := entry.Name()

		switch {
		case entry.IsDir():
			err = extractDir(source.Join(name), output.Join(name))
		case entry.Type().IsRegular():
			err = extractFile(source.WithFile(name), output.WithFile(name))
		default:
			return fmt.Errorf("unsupported file in archive: %q", source.WithFile(name))
		}

		if err != nil {
			return err
		}
	}

	if !created {
		return nil
	}

	if err := output.Chmod(info.Mode().Perm()); err != nil {
		return err
	}

	return os.Chtimes(output.Path(), info.ModTime(), info.ModTime())
}

// extractFile publishes a regular archive member with the permissions and timestamp go-getter extracted.
func extractFile(source, target file.File) error {
	target, _, err := extractTarget(target)
	if err != nil {
		return err
	}

	info, err := source.Info()
	if err != nil {
		return err
	}

	const ownerRead = 0o400

	if info.Mode().Perm()&ownerRead == 0 {
		if err := source.Chmod(info.Mode().Perm() | ownerRead); err != nil {
			return err
		}
	}

	return source.CopyAtomic(target, func(tmp file.File) error {
		if err := tmp.Chmod(info.Mode().Perm()); err != nil {
			return err
		}

		return os.Chtimes(tmp.Path(), info.ModTime(), info.ModTime())
	})
}

// extractTarget resolves an existing destination link and returns nil info for a new file.
func extractTarget(target file.File) (file.File, fs.FileInfo, error) {
	info, err := os.Lstat(target.Path())
	if errors.Is(err, fs.ErrNotExist) {
		return target, nil, nil
	}

	if err != nil {
		return "", nil, err
	}

	if info.Mode()&fs.ModeSymlink != 0 {
		path, err := filepath.EvalSymlinks(target.Path())
		if err != nil {
			return "", nil, fmt.Errorf("resolving output file %q: %w", target, err)
		}

		target = file.New(path)

		info, err = target.Info()
		if err != nil {
			return "", nil, err
		}
	}

	if !info.Mode().IsRegular() {
		return "", nil, fmt.Errorf("not a regular output file: %q", target)
	}

	return target, info, nil
}
