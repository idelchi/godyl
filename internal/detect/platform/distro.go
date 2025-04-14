package platform

import (
	"fmt"
	"strings"
)

// Distribution represents a Linux distribution configuration.
type Distribution struct {
	// Type is the canonical distribution name (e.g., debian, ubuntu).
	Type string

	// Raw contains the original string that was parsed.
	Raw string
}

// DistroInfo defines a Linux distribution's characteristics.
// Includes the canonical type name and known aliases.
type DistroInfo struct {
	Type    string
	Aliases []string
}

// Supported returns the list of supported Linux distributions.
// Includes major distributions like Debian, Ubuntu, CentOS, and their aliases.
func (DistroInfo) Supported() []DistroInfo {
	return []DistroInfo{
		{
			Type: "debian",
		},
		{
			Type: "ubuntu",
		},
		{
			Type: "centos",
		},
		{
			Type:    "redhat",
			Aliases: []string{"rhel"},
		},
		{
			Type: "arch",
		},
		{
			Type: "alpine",
		},
		{
			Type:    "raspbian",
			Aliases: []string{"raspberry"},
		},
	}
}

// Parse extracts distribution information from a string identifier.
// Matches against known distribution types and aliases, setting type
// and raw values. Returns an error if parsing fails.
func (d *Distribution) Parse(name string) error {
	name = strings.ToLower(name)

	info := DistroInfo{}

	for _, info := range info.Supported() {
		for _, alias := range append([]string{info.Type}, info.Aliases...) {
			if strings.Contains(name, alias) {
				d.Type = info.Type
				d.Raw = alias

				return nil
			}
		}
	}

	return fmt.Errorf("unable to parse distribution from name: %q", name)
}

// IsUnset checks if the distribution type is empty.
func (d Distribution) IsUnset() bool {
	return d.Type == ""
}

// Is checks for exact distribution match including raw string.
// Returns true only if both distributions are set and identical.
func (d Distribution) Is(other Distribution) bool {
	return other.Raw == d.Raw && !d.IsUnset() && !other.IsUnset()
}

// IsCompatibleWith checks if this distribution can run binaries built for another.
// Currently checks only for exact type matches between distributions.
func (d Distribution) IsCompatibleWith(other Distribution) bool {
	if d.IsUnset() || other.IsUnset() {
		return false
	}

	return d.Type == other.Type
}

// String returns the canonical name of the distribution.
func (d Distribution) String() string {
	return d.Type
}

// Default returns an empty Distribution configuration.
// Used as a fallback when distribution detection fails.
func Default() Distribution {
	return Distribution{}
}
