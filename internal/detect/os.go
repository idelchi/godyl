package detect

import (
	"fmt"

	"github.com/shirou/gopsutil/host"

	"github.com/idelchi/godyl/internal/detect/platform"
)

// Detect gathers information about the current platform, such as the operating system, architecture,
// distribution, library, and file extension, and populates the Platform struct accordingly.
func (p *Platform) Detect() error {
	var (
		operatingSystem platform.OS
		architecture    platform.Architecture
		library         platform.Library
		distro          platform.Distribution
		extension       platform.Extension
	)

	info, err := host.Info()
	if err != nil {
		return fmt.Errorf("getting host information: %w", err)
	}

	// Determine the OS from runtime information
	if err := operatingSystem.Parse(info.OS); err != nil {
		return fmt.Errorf("parsing OS: %w", err)
	}

	// Determine the Linux distribution from system information
	distro.Parse(info.Platform) //nolint:errcheck 	// Ignore error as it's not critical

	// Set the default library based on the OS and distribution
	library = library.Default(operatingSystem, distro)

	// Determine the architecture from the system's kernel architecture
	if err := architecture.Parse(info.KernelArch); err != nil {
		return fmt.Errorf("parsing architecture: %w", err)
	}

	if architecture.Is64Bit() && operatingSystem.Type == "linux" {
		is32Bit, err := platform.Is32Bit()
		if err == nil && is32Bit {
			architecture.To32BitUserLand()
		}
	}

	// Populate the Platform struct with the detected values
	*p = Platform{
		OS:           operatingSystem,
		Distribution: distro,
		Architecture: architecture,
		Library:      library,
		Extension:    extension.Default(operatingSystem),
	}

	return nil
}
