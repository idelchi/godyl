//go:build linux

package detect

import (
	"runtime"

	"github.com/idelchi/godyl/internal/detect/platform"
	"github.com/zcalusic/sysinfo"
)

// Detect gathers information about the current platform, such as the operating system, architecture,
// distribution, library, and file extension, and populates the Platform struct accordingly.
func (p *Platform) Detect() error {
	var os platform.OS
	var arch platform.Architecture
	var library platform.Library
	var distro platform.Distribution
	var extension platform.Extension

	var si sysinfo.SysInfo

	// Get system information
	si.GetSysInfo()

	// Determine the OS from runtime information
	if err := os.Parse(runtime.GOOS); err != nil {
		return err
	}

	// Determine the Linux distribution from system information
	distro.Parse(si.OS.Vendor)

	// Set the default library based on the OS and distribution
	library = library.Default(os, distro)

	// Determine the architecture from the system's kernel architecture
	if err := arch.Parse(runtime.GOARCH); err != nil {
		return err
	}

	// Populate the Platform struct with the detected values
	*p = Platform{
		OS:           os,
		Distribution: distro,
		Architecture: arch,
		Library:      library,
		Extension:    extension.Default(os),
	}

	return nil
}
