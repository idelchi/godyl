package platform

import (
	"fmt"
	"strings"

	"github.com/goccy/go-yaml/ast"

	"github.com/idelchi/godyl/pkg/unmarshal"
)

// Distribution represents a Linux distribution configuration.
type Distribution struct {
	Name string `single:"true"`

	// Type is the canonical distribution name (e.g., debian, ubuntu).
	canonical string

	// Raw contains the original string that was parsed.
	alias string
}

func (d *Distribution) IsNil() bool {
	return d.Name == ""
}

func (d *Distribution) UnmarshalYAML(node ast.Node) error {
	type raw Distribution

	if err := unmarshal.SingleStringOrStruct(node, (*raw)(d)); err != nil {
		return fmt.Errorf("unmarshaling distribution: %w", err)
	}

	return nil
}

func (d Distribution) MarshalYAML() (any, error) {
	return d.Name, nil
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
func (d *Distribution) ParseFrom(name string, comparisons ...func(string, string) bool) error {
	if len(comparisons) == 0 {
		comparisons = append(comparisons, strings.Contains)
	}

	lower := strings.ToLower(name)

	info := DistroInfo{}

	for _, info := range info.Supported() {
		for _, alias := range append([]string{info.Type}, info.Aliases...) {
			for _, compare := range comparisons {
				if compare(lower, alias) {
					d.Name = name
					d.canonical = info.Type
					d.alias = alias

					return nil
				}
			}
		}
	}

	return fmt.Errorf("unable to parse distribution from name: %q", name)
}

// Parse extracts operating system information from a string identifier.
func (d *Distribution) Parse() error {
	return d.ParseFrom(d.Name, strings.EqualFold, strings.Contains)
}

// IsUnset checks if the distribution type is empty.
func (d Distribution) IsUnset() bool {
	return d.canonical == ""
}

// Is checks for exact distribution match including raw string.
// Returns true only if both distributions are set and identical.
func (d Distribution) Is(other Distribution) bool {
	return other.alias == d.alias && !d.IsUnset() && !other.IsUnset()
}

// IsCompatibleWith checks if this distribution can run binaries built for another.
// Currently checks only for exact type matches between distributions.
func (d Distribution) IsCompatibleWith(other Distribution) bool {
	if d.IsUnset() || other.IsUnset() {
		return false
	}

	return d.canonical == other.canonical
}

// String returns the canonical name of the distribution.
func (d Distribution) String() string {
	return d.canonical
}

// Default returns an empty Distribution configuration.
// Used as a fallback when distribution detection fails.
func Default() Distribution {
	return Distribution{}
}
