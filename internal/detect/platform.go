package detect

import (
	"github.com/idelchi/godyl/internal/detect/platform"
	"github.com/idelchi/godyl/pkg/utils"
)

// Platform encapsulates system-specific characteristics and capabilities.
type Platform struct {
	// OS identifies the operating system (e.g., Linux, Windows, macOS).
	OS platform.OS

	// Architecture defines the CPU architecture and version (e.g., x86_64, arm64).
	Architecture platform.Architecture

	// Library specifies the system's standard library (e.g., GNU, Musl, MSVC).
	Library platform.Library

	// Extension defines platform-specific executable file extensions.
	Extension platform.Extension

	// Distribution identifies the Linux distribution when applicable.
	Distribution platform.Distribution
}

// Parse extracts platform information from a string identifier.
// Attempts to determine OS, architecture, and library details from the input.
func (p *Platform) Parse(name string) {
	p.OS.Parse(name)
	p.Architecture.Parse(name)
	p.Library.Parse(name)
	p.Extension.Default(p.OS)
}

// Merge combines two Platform configurations.
// Copies non-zero values from the other Platform into this one,
// preserving existing values when they are already set.
func (p *Platform) Merge(other Platform) {
	utils.SetIfZeroValue(&p.OS, other.OS)
	utils.SetIfZeroValue(&p.Architecture.Type, other.Architecture.Type)
	utils.SetIfZeroValue(&p.Architecture.Version, other.Architecture.Version)
	utils.SetIfZeroValue(&p.Architecture.Raw, other.Architecture.Raw)
	utils.SetIfZeroValue(&p.Library, other.Library)
	utils.SetIfZeroValue(&p.Extension, other.Extension)
	utils.SetIfZeroValue(&p.Distribution, other.Distribution)
}

// ToMap converts the Platform configuration into a map for templating.
// Includes derived values like architecture type, version, and capability flags.
func (p Platform) ToMap() map[string]any {
	platformMap := make(map[string]any)
	platformMap["OS"] = p.OS.String()
	platformMap["ARCH"] = p.Architecture.Type
	platformMap["ARCH_VERSION"] = p.Architecture.Version
	platformMap["ARCH_LONG"] = p.Architecture.String()
	platformMap["IS_ARM"] = p.Architecture.IsARM()
	platformMap["IS_X86"] = p.Architecture.IsX86()
	platformMap["LIBRARY"] = p.Library.String()
	platformMap["EXTENSION"] = p.Extension.String()
	platformMap["DISTRIBUTION"] = p.Distribution.String()

	return platformMap
}
