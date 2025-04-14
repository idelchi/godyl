package platform

import (
	"fmt"
	"strings"
)

// Library represents a system's standard library or ABI configuration.
type Library struct {
	// Type is the canonical library name (e.g., gnu, musl, msvc).
	Type string

	// Raw contains the original string that was parsed.
	Raw string
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
func (l *Library) Parse(name string) error {
	name = strings.ToLower(name)

	info := LibraryInfo{}

	for _, info := range info.Supported() {
		for _, alias := range append([]string{info.Type}, info.Aliases...) {
			if strings.Contains(name, alias) {
				l.Type = info.Type
				l.Raw = alias

				return nil
			}
		}
	}

	return fmt.Errorf("unable to parse library from name: %q", name)
}

// IsUnset checks if the library type is empty.
func (l Library) IsUnset() bool {
	return l.Type == ""
}

// IsSet checks if the library type has been configured.
func (l Library) IsSet() bool {
	return l.Type != ""
}

// Is checks for exact library match including raw string.
// Returns true only if both libraries are set and identical.
func (l Library) Is(other Library) bool {
	return other.Raw == l.Raw && !l.IsUnset() && !other.IsUnset()
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
	if compatible, exists := compatibilityMatrix[l.Type]; exists {
		return compatible[other.Type]
	}

	return false
}

// String returns the canonical name of the library.
func (l Library) String() string {
	return l.Type
}

// Default determines the standard library for a platform.
// Uses OS and distribution information to select the appropriate
// system library (e.g., GNU for Linux, MSVC for Windows).
func (l *Library) Default(os OS, distro Distribution) Library {
	switch os.Type {
	case "windows":
		return Library{Type: "msvc", Raw: "msvc"}
	case "android":
		return Library{Type: "android", Raw: "android"}
	case "linux":
		switch distro.Type {
		case "alpine":
			return Library{Type: "musl", Raw: "musl"}
		default:
			return Library{Type: "gnu", Raw: "gnu"}
		}
	default:
		return Library{}
	}
}
