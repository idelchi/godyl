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

// ErrParse is returned when a propery cannot be parsed.
var ErrParse = errors.New("unable to parse")

// Architecture-related constants.
const (
	armString  = "arm"
	armelValue = 5
	armhfValue = 7
	bit32Value = 32
)

// Architecture represents a CPU architecture with its type, version, and raw string.
type Architecture struct {
	Type            string
	Version         int
	Raw             string // Original parsed architecture value
	Is32BitUserLand bool
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

// Parse attempts to parse the architecture from the given name string.
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

// IsUnset returns true if the architecture type is not set.
func (a *Architecture) IsUnset() bool {
	return a.Type == ""
}

// Is checks if this architecture is exactly the same as another.
func (a *Architecture) Is(other Architecture) bool {
	return other.Raw == a.Raw && !a.IsUnset() && !other.IsUnset()
}

// IsCompatibleWith checks if this architecture is compatible with another.
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

// String returns a string representation of the architecture.
func (a *Architecture) String() string {
	if a.Version != 0 {
		if a.Type == armString {
			return fmt.Sprintf("armv%d", a.Version)
		}

		return fmt.Sprintf("%sv%d", a.Type, a.Version)
	}

	return a.Type
}

// To32BitUserLand converts 64-bit architecture to 32-bit equivalent.
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

// Is64Bit returns true if the architecture is 64-bit.
func (a *Architecture) Is64Bit() bool {
	return strings.Contains(a.Type, "64")
}

// Is32Bit returns true if the system is running in 32-bit mode.
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
