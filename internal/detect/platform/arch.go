package platform

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/sys/cpu"
)

// Architecture represents a CPU architecture with its type, version, and raw string.
type Architecture struct {
	Type    string
	Version int
	Raw     string // Original parsed architecture value
}

// ArchInfo holds information about an architecture type, including aliases and a parse function.
type ArchInfo struct {
	Type    string
	Aliases []string
	Parse   func(string) (int, error)
}

// Supported returns a slice of supported architecture information.
func (ArchInfo) Supported() []ArchInfo {
	return []ArchInfo{
		{
			Type:    "amd64",
			Aliases: []string{"amd64", "x86_64", "x64", "win64"},
		},
		{
			Type:    "386",
			Aliases: []string{"amd32", "x86", "i386", "i686", "386", "win32"},
		},
		{
			Type:    "arm64",
			Aliases: []string{"arm64", "aarch64"},
		},
		{
			Type:    "arm",
			Aliases: []string{"armv7", "armv6", "armv5", "armel", "armhf", "arm"},
			Parse: func(s string) (int, error) {
				switch s {
				case "armel":
					return 5, nil
				case "armhf":
					return 7, nil // (or 6)
				}

				re := regexp.MustCompile(`armv(\d+)`)
				match := re.FindStringSubmatch(s)
				if len(match) > 1 {
					return strconv.Atoi(match[1])
				}

				return getGOARM(), nil
			},
		},
	}
}

// Parse attempts to parse the architecture from the given name string.
func (a *Architecture) Parse(name string) error {
	name = strings.ToLower(name)

	info := ArchInfo{}

	for _, info := range info.Supported() {
		for _, alias := range info.Aliases {
			if strings.Contains(name, alias) {
				a.Type = info.Type
				a.Raw = alias
				if info.Parse != nil {
					version, err := info.Parse(alias)
					if err != nil {
						return err
					}
					a.Version = version
				}
				return nil
			}
		}
	}

	return fmt.Errorf("unable to parse architecture from name: %s", name)
}

// IsUnset returns true if the architecture type is not set.
func (a Architecture) IsUnset() bool {
	return a.Type == ""
}

// Is checks if this architecture is exactly the same as another.
func (a Architecture) Is(other Architecture) bool {
	return other.Raw == a.Raw && !a.IsUnset() && !other.IsUnset()
}

// IsCompatibleWith checks if this architecture is compatible with another.
func (a Architecture) IsCompatibleWith(other Architecture) bool {
	if a.IsUnset() || other.IsUnset() {
		return false
	}

	if a.Is(other) {
		return true
	}

	if a.Type != other.Type {
		return false
	}

	if a.Version != 0 && other.Version != 0 {
		return a.Version >= other.Version
	}
	return true
}

// String returns a string representation of the architecture.
func (a Architecture) String() string {
	if a.Version != 0 {
		if a.Type == "arm32" {
			return fmt.Sprintf("armv%d", a.Version)
		}

		return fmt.Sprintf("%sv%d", a.Type, a.Version)
	}
	return a.Type
}

func getGOARM() int {
	// Default to GOARM=5 if no special features are detected
	version := 5

	// Check for ARM CPU features using x/sys/cpu package
	if cpu.ARM.HasVFPv3 {
		version = 7 // ARMv7 with VFPv3 support
	} else if cpu.ARM.HasVFP {
		version = 6 // ARMv6 with VFP support
	}

	return version
}
