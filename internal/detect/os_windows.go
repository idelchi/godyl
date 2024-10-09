package detect

import (
	"runtime"

	"github.com/idelchi/godyl/internal/detect/platform"
)

func (p *Platform) Detect() error {
	var os platform.OS
	var arch platform.Architecture
	var library platform.Library = platform.MSVC
	var distro platform.Distribution
	var extension platform.Extension

	if err := os.From(runtime.GOOS); err != nil {
		return err
	}

	library = library.Default(os, distro)

	if err := arch.From(runtime.GOARCH, distro); err != nil {
		return err
	}

	*p = Platform{
		OS:           os,
		Distribution: distro,
		Architecture: arch,
		Library:      library,
		Extension:    extension.Default(os),
	}

	return nil
}
