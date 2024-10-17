package platform

import (
	"errors"
	"fmt"

	"github.com/idelchi/godyl/pkg/utils"
)

// OS represents an operating system type.
type OS string

// Predefined OS values.
const (
	Linux   OS = "linux"   // Linux operating system
	Darwin  OS = "darwin"  // Darwin operating system
	Windows OS = "windows" // Windows operating system
	FreeBSD OS = "freebsd" // FreeBSD operating system
	Android OS = "android" // Android operating system
	NetBSD  OS = "netbsd"  // NetBSD operating system
	OpenBSD OS = "openbsd" // OpenBSD operating system
)

// Available returns a slice of all supported operating systems.
func (o *OS) Available() []OS {
	return []OS{Linux, Darwin, Windows, FreeBSD, Android, NetBSD, OpenBSD}
}

// ErrNotFound is returned when no matching OS is found.
var ErrNotFound = errors.New("match not found")

// From sets the OS based on the provided string, if it matches any available OS or its aliases.
func (o *OS) From(operatingSystem string) error {
	for _, os := range o.Available() {
		if os.IsCompatibleWith(operatingSystem) {
			*o = os
			return nil
		}
	}
	return fmt.Errorf("%w: os %q", ErrNotFound, operatingSystem)
}

// CompatibleWith returns a list of strings that are compatible with the current OS.
func (o OS) CompatibleWith() []string {
	switch o {
	case Linux:
		return []string{"linux"}
	case Darwin:
		return []string{"darwin", "macos", "mac", "osx"}
	case Windows:
		return []string{"windows", "win"}
	case FreeBSD:
		return []string{"freebsd"}
	case Android:
		return []string{"android"}
	case NetBSD:
		return []string{"netbsd"}
	case OpenBSD:
		return []string{"openbsd"}
	default:
		return []string{}
	}
}

// Name returns the name of the OS as a string.
func (o OS) Name() string {
	return string(o)
}

// String returns the OS as a string.
func (o OS) String() string {
	return o.Name()
}

// IsCompatibleWith checks if the provided OS string matches any compatible OS aliases.
func (o OS) IsCompatibleWith(os string) bool {
	for _, compatible := range o.CompatibleWith() {
		if utils.EqualLower(os, compatible) {
			return true
		}
	}
	return false
}

// Parse attempts to parse a string and set the OS accordingly, based on its name or aliases.
func (o *OS) Parse(name string) error {
	for _, os := range o.Available() {
		if utils.ContainsLower(name, os.Name()) {
			*o = os
			return nil
		}
	}

	for _, os := range o.Available() {
		for _, alias := range os.CompatibleWith() {
			if utils.ContainsLower(name, alias) {
				*o = os
				return nil
			}
		}
	}

	// Return an error if no match is found
	return fmt.Errorf("unable to parse OS from name: %s", name)
}

// Default returns the default OS, which is Linux.
func (o *OS) Default() OS {
	return Linux
}
