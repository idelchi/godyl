//go:build !linux

package detect

import (
	"fmt"
	"runtime"

	"github.com/idelchi/godyl/internal/detect/platform"
)

// Detect gathers information about the current platform, including the operating system, architecture,
// library, and file extension, and populates the Platform struct accordingly for Windows and macOS.
func (p *Platform) Detect() error {
	var os platform.OS
	var arch platform.Architecture
	var library platform.Library
	var distro platform.Distribution
	var extension platform.Extension

	// Determine the OS from runtime information
	if err := os.Parse(runtime.GOOS); err != nil {
		return err
	}

	// Set the default library based on the OS (distribution is irrelevant for Windows/macOS)
	library = library.Default(os, distro)

	// Determine the architecture from runtime information
	if err := arch.Parse(runtime.GOARCH); err != nil {
		return err
	}

	if arch.Raw == "arm" {
		arch.Version = platform.InferGoArmVersion()
		arch.Raw = fmt.Sprintf("%sv%d", arch.Type, arch.Version)
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
