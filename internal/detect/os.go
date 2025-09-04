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

	operatingSystem = platform.OS{
		Name: info.OS,
	}
	// Determine the OS from runtime information
	if err := operatingSystem.ParseFrom(info.OS); err != nil {
		return fmt.Errorf("parsing OS: %w", err)
	}

	// Determine the Linux distribution from system information
	distro.ParseFrom(info.Platform) //nolint:gosec,errcheck 	// Ignore error as it's not critical

	// Determine the architecture from the system's kernel architecture
	if err := architecture.ParseFrom(info.KernelArch); err != nil {
		return fmt.Errorf("parsing architecture: %w", err)
	}

	// Set the default library based on the OS and distribution
	library = library.Default(operatingSystem, distro)

	if architecture.Is64Bit() && operatingSystem.Type() == "linux" {
		is32Bit, err := platform.Is32Bit()
		if err == nil && is32Bit {
			architecture.To32BitUserLand()
		}
	}

	// Set the default file extension based on the OS
	extension.ParseFrom(operatingSystem)

	// Populate the Platform struct with the detected values
	*p = Platform{
		OS:           operatingSystem,
		Distribution: distro,
		Architecture: architecture,
		Library:      library,
		Extension:    extension,
	}

	return nil
}
