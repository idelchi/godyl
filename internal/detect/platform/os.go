package platform

import (
	"fmt"
	"strings"

	"github.com/goccy/go-yaml/ast"
	"github.com/idelchi/godyl/pkg/unmarshal"
)

// OS represents an operating system configuration.
type OS struct {
	Name string `single:"true"`

	canonical string
	alias     string
}

func (o *OS) IsNil() bool {
	return o.Name == ""
}

func (o *OS) UnmarshalYAML(node ast.Node) error {
	type raw OS

	if err := unmarshal.SingleStringOrStruct(node, (*raw)(o)); err != nil {
		return fmt.Errorf("unmarshaling os: %w", err)
	}

	return nil
}

func (o OS) MarshalYAML() (any, error) {
	return o.Name, nil
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

func (o *OS) ParseFrom(name string, comparisons ...func(string, string) bool) error {
	if len(comparisons) == 0 {
		comparisons = append(comparisons, strings.Contains)
	}

	lower := strings.ToLower(name)

	osInfo := OSInfo{}

	for _, info := range osInfo.Supported() {
		for _, alias := range append([]string{info.Type}, info.Aliases...) {
			for _, compare := range comparisons {
				if compare(lower, alias) {
					o.Name = name
					o.canonical = info.Type
					o.alias = alias

					return nil
				}
			}
		}
	}

	return fmt.Errorf("%w: OS from name: %s", ErrParse, name)
}

// Parse extracts operating system information from a string identifier.
func (o *OS) Parse() error {
	return o.ParseFrom(o.Name, strings.EqualFold)
}

func (o *OS) Type() string {
	return o.canonical
}

// IsUnset checks if the OS type is empty.
func (o *OS) IsUnset() bool {
	return o.canonical == ""
}

// Is checks for exact OS match including raw string.
// Returns true only if both OS configurations are set and identical.
func (o *OS) Is(other OS) bool {
	return other.alias == o.alias && !o.IsUnset() && !other.IsUnset()
}

// IsCompatibleWith checks if binaries built for this OS can run on another.
// Currently requires exact OS type matches (e.g., linux-linux, windows-windows).
func (o *OS) IsCompatibleWith(other OS) bool {
	if o.IsUnset() || other.IsUnset() {
		return false
	}

	return o.canonical == other.canonical
}

// String returns the canonical name of the operating system.
func (o OS) String() string {
	return o.canonical
}
