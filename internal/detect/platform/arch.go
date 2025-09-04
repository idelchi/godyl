package platform

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml/ast"

	"github.com/idelchi/godyl/pkg/unmarshal"
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
	Name string `single:"true"`

	canonical       string
	alias           string
	version         int
	is32BitUserLand bool
}

// IsNil returns true if the Architecture pointer is nil.
func (a *Architecture) IsNil() bool {
	return a.Name == ""
}

// UnmarshalYAML implements the yaml.Unmarshaler interface for Architecture.
func (a *Architecture) UnmarshalYAML(node ast.Node) error {
	type raw Architecture

	if err := unmarshal.SingleStringOrStruct(node, (*raw)(a)); err != nil {
		return fmt.Errorf("unmarshaling architecture: %w", err)
	}

	return nil
}

// MarshalYAML implements the yaml.Marshaler interface for Architecture.
func (a *Architecture) MarshalYAML() (any, error) {
	return a.Name, nil
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

// ParseFrom extracts architecture information from a string identifier.
// Matches against known architecture types and aliases, setting type,
// version, and raw values. Returns an error if parsing fails.
//
//nolint:gocognit // Complex function - refactoring into smaller functions is a separate improvement task
func (a *Architecture) ParseFrom(name string, comparisons ...func(string, string) bool) error {
	if len(comparisons) == 0 {
		comparisons = append(comparisons, strings.Contains)
	}

	lower := strings.ToLower(name)

	archInfo := ArchInfo{}

	for _, info := range archInfo.Supported() {
		for i, alias := range append([]string{info.Type}, info.Aliases...) {
			if info.Type == armString && i == 0 {
				// Skip the arm type since it's the default and will be checked last
				continue
			}

			for _, compare := range comparisons {
				if compare(lower, alias) {
					a.Name = name
					a.canonical = info.Type
					a.alias = alias

					if info.Parse != nil {
						version, err := info.Parse(alias)
						if err != nil {
							return err
						}

						a.version = version
					}

					return nil
				}
			}
		}
	}

	return fmt.Errorf("%w: architecture from name: %s", ErrParse, name)
}

// Parse extracts operating system information from a string identifier.
func (a *Architecture) Parse() error {
	return a.ParseFrom(a.Name, strings.EqualFold)
}

// IsUnset checks if the architecture type is empty.
func (a *Architecture) IsUnset() bool {
	return a.canonical == ""
}

// IsSet checks if the architecture type has been configured.
func (a *Architecture) IsSet() bool {
	return a.canonical != ""
}

// Is checks for exact architecture match including raw string.
// Returns true only if both architectures are set and identical.
func (a *Architecture) Is(other Architecture) bool {
	return other.alias == a.alias && !a.IsUnset() && !other.IsUnset()
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

	if a.canonical != other.canonical {
		return false
	}

	if a.version != 0 && other.version != 0 {
		return a.version >= other.version
	}

	return true
}

// String returns the canonical string representation of the architecture.
// Includes version information for ARM architectures (e.g., "armv7").
func (a *Architecture) String() string {
	if a.version != 0 {
		if a.canonical == armString {
			return fmt.Sprintf("armv%d", a.version)
		}

		return fmt.Sprintf("%sv%d", a.canonical, a.version)
	}

	return a.canonical
}

// To32BitUserLand converts a 64-bit architecture to its 32-bit equivalent.
// Handles amd64->386 and arm64->armv7 conversions for 32-bit userland support.
func (a *Architecture) To32BitUserLand() {
	a.is32BitUserLand = true

	switch a.canonical {
	case "amd64":
		a.canonical = "386"
	case "arm64":
		a.canonical = armString
		a.version = armhfValue
		a.alias = "armv7"
	}
}

// Type returns the canonical architecture type name.
func (a *Architecture) Type() string {
	return a.canonical
}

// Version returns the architecture version number.
func (a *Architecture) Version() int {
	return a.version
}

// Is64Bit checks if the architecture is 64-bit capable.
// Returns true for amd64 and arm64 architectures.
func (a *Architecture) Is64Bit() bool {
	return strings.Contains(a.canonical, "64")
}

// IsX86 checks if the architecture is x86-based.
// Returns true for both 32-bit (386) and 64-bit (amd64) variants.
func (a *Architecture) IsX86() bool {
	return a.canonical == "amd64" || a.canonical == "386"
}

// IsARM checks if the architecture is ARM-based.
// Returns true for both 32-bit (arm) and 64-bit (arm64) variants.
func (a *Architecture) IsARM() bool {
	return a.canonical == armString || a.canonical == "arm64"
}

// Is32Bit detects if the system is running in 32-bit mode.
// Uses getconf to determine the system's bit width.
func Is32Bit() (bool, error) {
	cmd := exec.CommandContext(context.Background(), "getconf", "LONG_BIT")

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
