package platform

import (
	"fmt"
	"strings"
)

// OS represents an operating system configuration.
type OS struct {
	// Type is the canonical OS name (e.g., linux, windows, darwin).
	Type string

	// Raw contains the original string that was parsed.
	Raw string
}

// OSInfo defines an operating system's characteristics.
// Includes the canonical type name and known aliases.
type OSInfo struct {
	Type    string
	Aliases []string
}

// Supported returns the list of supported operating systems.
// Includes major operating systems like Linux, macOS, Windows,
// and various BSD variants with their common aliases.
func (OSInfo) Supported() []OSInfo {
	return []OSInfo{
		{
			Type: "linux",
		},
		{
			Type:    "darwin",
			Aliases: []string{"macos", "mac", "osx"},
		},
		{
			Type:    "windows",
			Aliases: []string{"win"},
		},
		{
			Type: "freebsd",
		},
		{
			Type: "android",
		},
		{
			Type: "netbsd",
		},
		{
			Type: "openbsd",
		},
	}
}

// Parse extracts operating system information from a string identifier.
// Matches against known OS types and aliases, setting type and raw values.
// Returns an error if parsing fails.
func (o *OS) Parse(name string) error {
	name = strings.ToLower(name)

	osInfo := OSInfo{}

	for _, info := range osInfo.Supported() {
		for _, alias := range append([]string{info.Type}, info.Aliases...) {
			if strings.Contains(name, alias) {
				o.Type = info.Type
				o.Raw = alias

				return nil
			}
		}
	}

	return fmt.Errorf("%w: OS from name: %s", ErrParse, name)
}

// IsUnset checks if the OS type is empty.
func (o *OS) IsUnset() bool {
	return o.Type == ""
}

// Is checks for exact OS match including raw string.
// Returns true only if both OS configurations are set and identical.
func (o *OS) Is(other OS) bool {
	return other.Raw == o.Raw && !o.IsUnset() && !other.IsUnset()
}

// IsCompatibleWith checks if binaries built for this OS can run on another.
// Currently requires exact OS type matches (e.g., linux-linux, windows-windows).
func (o *OS) IsCompatibleWith(other OS) bool {
	if o.IsUnset() || other.IsUnset() {
		return false
	}

	return o.Type == other.Type
}

// String returns the canonical name of the operating system.
func (o OS) String() string {
	return o.Type
}
