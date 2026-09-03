package platform

import (
	"fmt"
	"strings"

	"github.com/goccy/go-yaml/ast"

	"github.com/idelchi/godyl/pkg/unmarshal"
)

const (
	// libAndroid is the canonical Android ABI name.
	libAndroid = osAndroid
	// libGNU is the canonical GNU C library name.
	libGNU = "gnu"
	// libMSVC is the canonical Microsoft Visual C++ runtime name.
	libMSVC = "msvc"
	// libMusl is the canonical musl C library name.
	libMusl = "musl"
	// libSystem is the canonical macOS system library name.
	libSystem = "libSystem"
)

// Library represents a system's standard library or ABI configuration.
type Library struct {
	// Name retains the unparsed library name supplied by the user or detector.
	Name string `single:"true"`
	// canonical is the normalized library name, such as gnu, musl, or msvc.
	canonical string

	// alias is the exact supported spelling that matched Name.
	alias string
}

// IsNil reports whether the unparsed library name is empty.
func (l *Library) IsNil() bool {
	return l.Name == ""
}

// UnmarshalYAML implements the yaml.Unmarshaler interface for Library.
func (l *Library) UnmarshalYAML(node ast.Node) error {
	type raw Library

	if err := unmarshal.SingleStringOrStruct(node, (*raw)(l)); err != nil {
		return fmt.Errorf("unmarshaling library: %w", err)
	}

	return nil
}

// MarshalYAML implements the yaml.Marshaler interface for Library.
func (l *Library) MarshalYAML() (any, error) {
	return l.Name, nil
}

// LibraryInfo defines a system library's characteristics.
// Includes the canonical type name and known aliases.
type LibraryInfo struct {
	// Type is the canonical library name.
	Type string
	// Aliases contains additional recognized spellings.
	Aliases []string
}

// Supported returns the list of supported system libraries.
// Includes major libraries like GNU, Musl, MSVC, and their aliases.
func (LibraryInfo) Supported() []LibraryInfo {
	return []LibraryInfo{
		{
			Type:    libGNU,
			Aliases: []string{"glibc"},
		},
		{
			Type: libMusl,
		},
		{
			Type:    libMSVC,
			Aliases: []string{"visualcpp"},
		},
		{
			Type: libAndroid,
		},
		{
			Type: libSystem,
		},
	}
}

// ParseFrom extracts library information from a string identifier.
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

	return fmt.Errorf("%w: library from %q", ErrParse, name)
}

// Parse extracts library information from Name.
func (l *Library) Parse() error {
	return l.ParseFrom(l.Name, strings.EqualFold, strings.Contains)
}

// IsUnset checks if the library type is empty.
func (l *Library) IsUnset() bool {
	return l.canonical == ""
}

// IsSet checks if the library type has been configured.
func (l *Library) IsSet() bool {
	return l.canonical != ""
}

// Is checks for exact library match including raw string.
// Returns true only if both libraries are set and identical.
func (l *Library) Is(other Library) bool {
	return other.alias == l.alias && !l.IsUnset() && !other.IsUnset()
}

// compatibilityMatrix records which target libraries can run binaries built against each source library.
var compatibilityMatrix = map[string]map[string]bool{ //nolint:gochecknoglobals,lll // Compatibility matrix lookup table is appropriate as global
	libGNU: {
		libGNU:  true,
		libMusl: true,
	},
	libMusl: {
		libGNU:  true,
		libMusl: true,
	},
	libMSVC: {
		libMSVC: true,
		libGNU:  true,
	},
	libAndroid: {
		libAndroid: true,
	},
	libSystem: {
		libSystem: true,
		libGNU:    true,
		libMusl:   true,
	},
}

// IsCompatibleWith checks if binaries built against this library can run
// with another library. Uses a compatibility matrix to determine binary
// compatibility between different library implementations.
func (l *Library) IsCompatibleWith(other Library) bool {
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
func (l *Library) String() string {
	return l.canonical
}

// Default determines the standard library for a platform.
// Uses OS and distribution information to select the appropriate
// system library (e.g., GNU for Linux, MSVC for Windows).
func (l *Library) Default(os OS, distro Distribution) Library {
	switch os.Type() {
	case osWindows:
		return Library{Name: libMSVC, canonical: libMSVC, alias: libMSVC}
	case osAndroid:
		return Library{Name: libAndroid, canonical: libAndroid, alias: libAndroid}
	case osLinux:
		switch distro.canonical {
		case distroAlpine:
			return Library{Name: libMusl, canonical: libMusl, alias: libMusl}
		default:
			return Library{Name: libGNU, canonical: libGNU, alias: libGNU}
		}
	case osDarwin:
		return Library{Name: libSystem, canonical: libSystem, alias: libSystem}
	default:
		return Library{}
	}
}
