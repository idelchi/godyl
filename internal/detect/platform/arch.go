package platform

import (
	"fmt"

	"github.com/idelchi/godyl/pkg/compare"
)

type Type string

type Architecture struct {
	Type    Type
	Version string
}

func (a Architecture) Name() string {
	if a.Version != "" {
		return fmt.Sprintf("%sv%s", a.Type, a.Version)
	}
	return string(a.Type)
}

const (
	AMD64 Type = "amd64"
	AMD32 Type = "amd32"
	ARM64 Type = "arm64"
	ARM32 Type = "arm"
	MIPS  Type = "mips"
	PPC64 Type = "ppc64"
	RISCV Type = "riscv"
)

func (a Architecture) Default() Architecture {
	return Architecture{
		Type: AMD64,
	}
}

func (a Architecture) Supported() []Architecture {
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

func (a *Architecture) From(architecture string, distro Distribution) error {
	for _, arch := range a.Supported() {
		if compare.Lower(architecture, arch.Name()) {
			*a = arch

			return nil
		}
	}

	for _, arch := range a.Supported() {
		if arch.IsCompatibleWith(architecture, distro) {
			*a = arch

			return nil
		}
	}

	return fmt.Errorf("%w: architecture %q", ErrNotFound, architecture)
}

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
				return []string{"arm32", "armhf", "armv6", "armv6l", "arm"}
			}
			if distro == Debian {
				return []string{"arm32", "armv6", "armv6l", "arm"}
			}
			return []string{"arm32", "armhf", "armv6", "armv6l", "arm"}
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

func (a Architecture) String() string {
	return a.Name()
}

func (a Architecture) IsCompatibleWith(arch string, distro Distribution) bool {
	for _, compatible := range a.CompatibleWith(distro) {
		if arch == compatible {
			return true
		}
	}
	return false
}

func (a *Architecture) Parse(name string) error {
	for _, arch := range a.Supported() {
		if compare.ContainsLower(name, arch.Name()) {
			*a = arch
			return nil
		}
	}

	for _, arch := range a.Supported() {
		for _, alias := range arch.CompatibleWith("") {
			if compare.ContainsLower(name, alias) {
				*a = arch
				return nil
			}
		}
	}

	// Return an error if no match is found
	return fmt.Errorf("unable to parse architecture from name: %s", name)
}
