package detect

import (
	"runtime"

	"github.com/idelchi/godyl/internal/detect/platform"
	"github.com/zcalusic/sysinfo"
)

func (p *Platform) Detect() error {
	var os platform.OS
	var arch platform.Architecture
	var library platform.Library = platform.GNU
	var distro platform.Distribution
	var extension platform.Extension

	var si sysinfo.SysInfo

	si.GetSysInfo()

	if err := os.From(runtime.GOOS); err != nil {
		return err
	}

	if err := distro.From(si.OS.Vendor); err != nil {
		return err
	}

	library = library.Default(os, distro)

	if err := arch.From(si.Kernel.Architecture, distro); err != nil {
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
