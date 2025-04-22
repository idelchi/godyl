package platform

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// ErrParse indicates a failure to parse platform-specific information.
var ErrParse = errors.New("unable to parse")

// Architecture-related constants.
const (
	armString  = "arm"
	armelValue = 5
	armhfValue = 7
	bit32Value = 32
)

// Architecture represents a CPU architecture configuration.
// Tracks architecture type, version, raw string, and user-land bitness.
type Architecture struct {
	Type            string
	Raw             string
	Version         int
	Is32BitUserLand bool
}

// ArchInfo defines an architecture's characteristics and parsing rules.
// Includes the canonical type name, known aliases, and version parsing logic.
type ArchInfo struct {
	Parse   func(string) (int, error)
	Type    string
	Aliases []string
}

// Supported returns the list of supported CPU architectures.
// Includes x86, ARM architectures and their variants with parsing rules.
func (ArchInfo) Supported() []ArchInfo {
	return []ArchInfo{
		{
			Type:    "amd64",
			Aliases: []string{"x86_64", "x64", "win64"},
		},
		{
			Type:    "386",
			Aliases: []string{"amd32", "x86", "i386", "i686", "win32"},
		},
		{
			Type:    "arm64",
			Aliases: []string{"aarch64"},
		},
		{
			Type:    armString,
			Aliases: []string{"armv7", "armv6", "armv5", "armel", "armhf", armString},
			Parse: func(str string) (int, error) {
				switch str {
				case "armel", armString:
					return armelValue, nil
				case "armhf":
					return armhfValue, nil // (or 6)
				}

				re := regexp.MustCompile(`armv(\d+)`)
				match := re.FindStringSubmatch(str)
				if len(match) > 1 {
					return strconv.Atoi(match[1])
				}

				return armelValue, nil
			},
		},
	}
}

// Parse extracts architecture information from a string identifier.
// Matches against known architecture types and aliases, setting type,
// version, and raw values. Returns an error if parsing fails.
func (a *Architecture) Parse(name string) error {
	name = strings.ToLower(name)

	archInfo := ArchInfo{}

	for _, info := range archInfo.Supported() {
		for i, alias := range append([]string{info.Type}, info.Aliases...) {
			if info.Type == armString && i == 0 {
				// Skip the arm type since it's the default and will be checked last
				continue
			}

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

	return fmt.Errorf("%w: architecture from name: %s", ErrParse, name)
}

// IsUnset checks if the architecture type is empty.
func (a *Architecture) IsUnset() bool {
	return a.Type == ""
}

// IsSet checks if the architecture type has been configured.
func (a *Architecture) IsSet() bool {
	return a.Type != ""
}

// Is checks for exact architecture match including raw string.
// Returns true only if both architectures are set and identical.
func (a *Architecture) Is(other Architecture) bool {
	return other.Raw == a.Raw && !a.IsUnset() && !other.IsUnset()
}

// IsCompatibleWith checks if this architecture can run binaries built for another.
// Considers architecture type and version compatibility (e.g., armv7 can run armv5).
func (a *Architecture) IsCompatibleWith(other Architecture) bool {
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

// String returns the canonical string representation of the architecture.
// Includes version information for ARM architectures (e.g., "armv7").
func (a Architecture) String() string {
	if a.Version != 0 {
		if a.Type == armString {
			return fmt.Sprintf("armv%d", a.Version)
		}

		return fmt.Sprintf("%sv%d", a.Type, a.Version)
	}

	return a.Type
}

// To32BitUserLand converts a 64-bit architecture to its 32-bit equivalent.
// Handles amd64->386 and arm64->armv7 conversions for 32-bit userland support.
func (a *Architecture) To32BitUserLand() {
	a.Is32BitUserLand = true

	switch a.Type {
	case "amd64":
		a.Type = "386"
	case "arm64":
		a.Type = armString
		a.Version = armhfValue
		a.Raw = "armv7"
	}
}

// Is64Bit checks if the architecture is 64-bit capable.
// Returns true for amd64 and arm64 architectures.
func (a *Architecture) Is64Bit() bool {
	return strings.Contains(a.Type, "64")
}

// IsX86 checks if the architecture is x86-based.
// Returns true for both 32-bit (386) and 64-bit (amd64) variants.
func (a *Architecture) IsX86() bool {
	return a.Type == "amd64" || a.Type == "386"
}

// IsARM checks if the architecture is ARM-based.
// Returns true for both 32-bit (arm) and 64-bit (arm64) variants.
func (a *Architecture) IsARM() bool {
	return a.Type == armString || a.Type == "arm64"
}

// Is32Bit detects if the system is running in 32-bit mode.
// Uses getconf to determine the system's bit width.
func Is32Bit() (bool, error) {
	cmd := exec.Command("getconf", "LONG_BIT")

	var out bytes.Buffer

	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return false, fmt.Errorf("running getconf command: %w", err)
	}

	result := strings.TrimSpace(out.String())

	value, err := strconv.Atoi(result)
	if err != nil {
		return false, fmt.Errorf("parsing bit value: %w", err)
	}

	return value == bit32Value, nil
}
