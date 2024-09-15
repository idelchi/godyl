package platform

import (
	"fmt"
	"strings"

	"github.com/goccy/go-yaml/ast"
	"github.com/idelchi/godyl/pkg/unmarshal"
)

// Library represents a system's standard library or ABI configuration.
type Library struct {
	Name string
	// Type is the canonical library name (e.g., gnu, musl, msvc).
	canonical string

	// Raw contains the original string that was parsed.
	alias string
}

func (l *Library) IsNil() bool {
	return l.Name == ""
}

func (l *Library) UnmarshalYAML(node ast.Node) error {
	type raw Library

	if err := unmarshal.SingleStringOrStruct(node, (*raw)(l)); err != nil {
		return fmt.Errorf("unmarshaling os: %w", err)
	}

	return nil
}

func (l Library) MarshalYAML() (any, error) {
	return l.Name, nil
}

// LibraryInfo defines a system library's characteristics.
// Includes the canonical type name and known aliases.
type LibraryInfo struct {
	Type    string
	Aliases []string
}

// Supported returns the list of supported system libraries.
// Includes major libraries like GNU, Musl, MSVC, and their aliases.
func (LibraryInfo) Supported() []LibraryInfo {
	return []LibraryInfo{
		{
			Type:    "gnu",
			Aliases: []string{"glibc"},
		},
		{
			Type: "musl",
		},
		{
			Type:    "msvc",
			Aliases: []string{"visualcpp"},
		},
		{
			Type: "android",
		},
	}
}

// Parse extracts library information from a string identifier.
// Matches against known library types and aliases, setting type
// and raw values. Returns an error if parsing fails.
func (l *Library) ParseFrom(name string, comparisons ...func(string, string) bool) error {
	if len(comparisons) == 0 {
		comparisons = append(comparisons, strings.Contains)
	}

	lower := strings.ToLower(name)

	info := LibraryInfo{}

	for _, info := range info.Supported() {
		for _, alias := range append([]string{info.Type}, info.Aliases...) {
			for _, compare := range comparisons {
				if compare(lower, alias) {
					l.Name = name
					l.canonical = info.Type
					l.alias = alias

					return nil
				}
			}
		}
	}

	return fmt.Errorf("unable to parse library from name: %q", name)
}

// Parse extracts operating system information from a string identifier.
func (l *Library) Parse() error {
	return l.ParseFrom(l.Name, strings.EqualFold, strings.Contains)
}

// IsUnset checks if the library type is empty.
func (l Library) IsUnset() bool {
	return l.canonical == ""
}

// IsSet checks if the library type has been configured.
func (l Library) IsSet() bool {
	return l.canonical != ""
}

// Is checks for exact library match including raw string.
// Returns true only if both libraries are set and identical.
func (l Library) Is(other Library) bool {
	return other.alias == l.alias && !l.IsUnset() && !other.IsUnset()
}

var compatibilityMatrix = map[string]map[string]bool{
	"gnu": {
		"gnu":  true,
		"musl": true,
	},
	"musl": {
		"gnu":  true,
		"musl": true,
	},
	"msvc": {
		"msvc": true,
		"gnu":  true,
	},
	"android": {
		"android": true,
	},
}

// IsCompatibleWith checks if binaries built against this library can run
// with another library. Uses a compatibility matrix to determine binary
// compatibility between different library implementations.
func (l Library) IsCompatibleWith(other Library) bool {
	// Early return if either library is unset
	if l.IsUnset() || other.IsUnset() {
		return false
	}

	// Check if they're exactly the same library
	if l.Is(other) {
		return true
	}

	// Look up compatibility in the matrix
	if compatible, exists := compatibilityMatrix[l.canonical]; exists {
		return compatible[other.canonical]
	}

	return false
}

// String returns the canonical name of the library.
func (l Library) String() string {
	return l.canonical
}

// Default determines the standard library for a platform.
// Uses OS and distribution information to select the appropriate
// system library (e.g., GNU for Linux, MSVC for Windows).
func (l *Library) Default(os OS, distro Distribution) Library {
	switch os.Type() {
	case "windows":
		return Library{Name: "msvc", canonical: "msvc", alias: "msvc"}
	case "android":
		return Library{Name: "android", canonical: "android", alias: "android"}
	case "linux":
		switch distro.canonical {
		case "alpine":
			return Library{Name: "musl", canonical: "musl", alias: "musl"}
		default:
			return Library{Name: "gnu", canonical: "gnu", alias: "gnu"}
		}
	default:
		return Library{}
	}
}
