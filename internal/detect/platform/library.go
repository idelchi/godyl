package platform

import (
	"fmt"
	"strings"
)

// Library represents a system library or ABI (Application Binary Interface).
type Library struct {
	Type string
	Raw  string // Original parsed library value
}

// LibraryInfo holds information about a library type, including aliases.
type LibraryInfo struct {
	Type    string
	Aliases []string
}

// Supported returns a slice of supported library information.
func (LibraryInfo) Supported() []LibraryInfo {
	return []LibraryInfo{
		{
			Type:    "gnu",
			Aliases: []string{"gnu", "glibc"},
		},
		{
			Type:    "musl",
			Aliases: []string{"musl"},
		},
		{
			Type:    "msvc",
			Aliases: []string{"msvc", "visualcpp"},
		},
		{
			Type:    "android",
			Aliases: []string{"android"},
		},
	}
}

// Parse attempts to parse the library from the given name string.
func (l *Library) Parse(name string) error {
	name = strings.ToLower(name)

	info := LibraryInfo{}

	for _, info := range info.Supported() {
		for _, alias := range info.Aliases {
			if strings.Contains(name, alias) {
				l.Type = info.Type
				l.Raw = alias

				return nil
			}
		}
	}

	return fmt.Errorf("unable to parse library from name: %s", name)
}

// IsUnset returns true if the library type is not set.
func (l Library) IsUnset() bool {
	return l.Type == ""
}

// Is checks if this library is exactly the same as another.
func (l Library) Is(other Library) bool {
	return other.Raw == l.Raw && !l.IsUnset() && !other.IsUnset()
}

// IsCompatibleWith checks if this library is compatible with another.
func (l Library) IsCompatibleWith(other Library) bool {
	if l.IsUnset() || other.IsUnset() {
		return false
	}

	if l.Is(other) {
		return true
	}

	if l.Type == other.Type {
		return true
	}

	// Specific compatibility rules
	if l.Type == "gnu" && other.Type == "musl" {
		return true
	}
	if l.Type == "musl" && other.Type == "gnu" {
		return true
	}
	if l.Type == "msvc" && other.Type == "gnu" {
		return true
	}

	return false
}

// String returns a string representation of the library.
func (l Library) String() string {
	return l.Type
}

// Default returns the default Library for a given OS and Distribution.
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
