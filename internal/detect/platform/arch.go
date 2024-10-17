package platform

import (
	"fmt"

	"github.com/idelchi/godyl/pkg/utils"
)

// Type represents a CPU architecture type, such as "amd64" or "arm64".
type Type string

// Architecture defines a platform architecture, consisting of a type and an optional version.
type Architecture struct {
	Type    Type   // Type indicates the architecture type (e.g., amd64, arm64).
	Version string // Version specifies the version for architectures that have variants, such as arm.
}

// Name returns the full name of the architecture, including version if present.
func (a Architecture) Name() string {
	if a.Version != "" {
		return fmt.Sprintf("%sv%s", a.Type, a.Version)
	}
	return string(a.Type)
}

// Short returns a short identifier for certain architectures.
func (a Architecture) Short() string {
	switch a.Type {
	case AMD64, ARM64:
		return "amd"
	case AMD32, ARM32:
		return "arm"
	default:
		return string(a.Type)
	}
}

const (
	AMD64 Type = "amd64" // AMD64 represents the 64-bit x86 architecture.
	AMD32 Type = "amd32" // AMD32 represents the 32-bit x86 architecture.
	ARM64 Type = "arm64" // ARM64 represents the 64-bit ARM architecture.
	ARM32 Type = "arm"   // ARM32 represents the 32-bit ARM architecture.
	MIPS  Type = "mips"  // MIPS represents the MIPS architecture.
	PPC64 Type = "ppc64" // PPC64 represents the 64-bit PowerPC architecture.
	RISCV Type = "riscv" // RISCV represents the RISC-V architecture.
)

// Default returns the default Architecture, which is AMD64.
func (a Architecture) Default() Architecture {
	return Architecture{
		Type: AMD64,
	}
}

// Available returns a list of all supported architectures.
func (a Architecture) Available() []Architecture {
	return []Architecture{
		{AMD64, ""},
		{AMD32, ""},
		{ARM64, ""},
		{ARM32, "6"},
		{ARM32, "7"},
		{MIPS, ""},
		{PPC64, ""},
		{RISCV, ""},
	}
}

// From sets the architecture based on the provided name and distribution, if found.
func (a *Architecture) From(architecture string, distro Distribution) error {
	for _, arch := range a.Available() {
		if utils.EqualLower(architecture, arch.Name()) {
			*a = arch
			return nil
		}
	}

	for _, arch := range a.Available() {
		if arch.IsCompatibleWith(architecture, distro) {
			*a = arch
			return nil
		}
	}

	return fmt.Errorf("%w: architecture %q", ErrNotFound, architecture)
}

// CompatibleWith returns a list of architecture aliases that are compatible with the given distribution.
func (a Architecture) CompatibleWith(distro Distribution) []string {
	switch a.Type {
	case AMD64:
		return []string{"amd64", "x86_64", "x64", "win64"}
	case ARM64:
		return []string{"arm64", "aarch64"}
	case AMD32:
		return []string{"amd32", "x86", "i386", "i686", "win32", "386", "686"}
	case ARM32:
		switch a.Version {
		case "7":
			return []string{"arm32", "armv7", "armv7l", "armhf", "armv6", "armv6l", "arm"}
		case "6":
			if distro == Rasbian {
				return []string{"arm32", "armv6", "armv6l", "armhf", "arm"}
			}
			if distro == Debian {
				return []string{"arm32", "armv6", "armv6l", "arm"}
			}
			return []string{"arm32", "armhf", "armv6", "armv6l", "arm"}
		case "5":
			return []string{"arm32", "armv5", "armv5l", "armel", "arm"}
		default:
			return []string{"arm32", "armv7", "armv7l", "armhf", "armv6", "armv6l", "arm"}
		}
	case MIPS:
		return []string{"mips", "mipsel", "mips64", "mips64le"}
	case PPC64:
		return []string{"powerpc64", "ppc64", "ppc64le"}
	case RISCV:
		return []string{"riscv64"}
	default:
		return []string{}
	}
}

// String returns the architecture name as a string.
func (a Architecture) String() string {
	return a.Name()
}

// IsCompatibleWith checks if the provided architecture string matches any of the aliases for this architecture.
func (a Architecture) IsCompatibleWith(arch string, distro Distribution) bool {
	for _, compatible := range a.CompatibleWith(distro) {
		if arch == compatible {
			return true
		}
	}
	return false
}

// Parse attempts to parse a string and set the architecture accordingly, based on its name or aliases.
func (a *Architecture) Parse(name string) error {
	for _, arch := range a.Available() {
		if utils.ContainsLower(name, arch.Name()) {
			*a = arch

			return nil
		}
	}

	for _, arch := range a.Available() {
		// TODO(Idelchi): Why arch.CompatibleWith("")? Isn't the distro required?
		for _, alias := range arch.CompatibleWith("") {
			if utils.ContainsLower(name, alias) {
				*a = arch

				return nil
			}
		}
	}

	// Return an error if no match is found
	return fmt.Errorf("unable to parse architecture from name: %s", name)
}
