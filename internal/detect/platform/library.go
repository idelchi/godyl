package platform

import (
	"fmt"

	"github.com/idelchi/godyl/pkg/utils"
)

// Library represents a system library or ABI (Application Binary Interface) used by an operating system or platform.
type Library string

// Predefined Library values.
const (
	GNU        Library = "gnu"     // GNU library, typically used in Linux distributions.
	Musl       Library = "musl"    // Musl library, often used in lightweight Linux distributions like Alpine.
	MSVC       Library = "msvc"    // Microsoft Visual C++ (MSVC) library, used in Windows.
	LibAndroid Library = "android" // Android library, used in Android OS.
)

// Default returns the default Library for a given OS and Distribution.
func (l *Library) Default(os OS, distro Distribution) Library {
	switch os {
	case Windows:
		return MSVC
	case Android:
		return LibAndroid
	case Linux:
		switch distro {
		case Alpine:
			return Musl
		default:
			return GNU
		}
	default:
		return ""
	}
}

// Available returns a slice of all supported libraries.
func (l Library) Available() []Library {
	return []Library{Musl, GNU, MSVC, LibAndroid}
}

// From sets the Library based on the provided string, if it matches any available library.
func (l *Library) From(library string) error {
	for _, lib := range l.Available() {
		if utils.EqualLower(library, lib.Name()) {
			*l = lib
			return nil
		}
	}

	for _, lib := range l.Available() {
		if lib.IsCompatibleWith(library) {
			*l = lib
			return nil
		}
	}

	*l = "" // Reset to empty if no match is found
	return nil
}

// CompatibleWith returns a list of compatible library names for the given Library.
func (l Library) CompatibleWith() []string {
	switch l {
	case GNU:
		return []string{"gnu", "musl"}
	case Musl:
		return []string{"musl", "gnu"}
	case MSVC:
		return []string{"msvc", "gnu"}
	case LibAndroid:
		return []string{"android"}
	}
	return nil
}

// IsCompatibleWith checks if the provided library name is compatible with the current Library.
func (l Library) IsCompatibleWith(lib string) bool {
	for _, compatible := range l.CompatibleWith() {
		if lib == compatible {
			return true
		}
	}
	return false
}

// Name returns the name of the Library as a string.
func (l Library) Name() string {
	return string(l)
}

// String returns the Library as a string.
func (l Library) String() string {
	return l.Name()
}

// Parse attempts to parse a string and set the Library accordingly, based on its name or aliases.
func (l *Library) Parse(name string) error {
	for _, library := range l.Available() {
		if utils.ContainsLower(name, library.Name()) {
			*l = library
			return nil
		}
	}

	for _, library := range l.Available() {
		for _, alias := range library.CompatibleWith() {
			if utils.ContainsLower(name, alias) {
				*l = library
				return nil
			}
		}
	}

	// Return an error if no match is found
	return fmt.Errorf("unable to parse library from name: %s", name)
}
