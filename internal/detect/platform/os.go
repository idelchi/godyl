package platform

import (
	"errors"
	"fmt"

	"github.com/idelchi/godyl/pkg/compare"
)

type OS string

const (
	Linux   OS = "linux"
	MacOS   OS = "macos"
	Windows OS = "windows"
	FreeBSD OS = "freebsd"
	Android OS = "android"
	NetBSD  OS = "netbsd"
	OpenBSD OS = "openbsd"
)

func (o *OS) Supported() []OS {
	return []OS{Linux, MacOS, Windows, FreeBSD, Android, NetBSD, OpenBSD}
}

var ErrNotFound = errors.New("match not found")

func (o *OS) From(operatingSystem string) error {
	for _, os := range o.Supported() {
		if os.IsCompatibleWith(operatingSystem) {
			*o = os

			return nil
		}
	}

	return fmt.Errorf("%w: os %q", ErrNotFound, operatingSystem)
}

func (o OS) CompatibleWith() []string {
	switch o {
	case Linux:
		return []string{"linux"}
	case MacOS:
		return []string{"macos", "mac", "osx", "darwin"}
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

func (o OS) Name() string {
	return string(o)
}

func (o OS) String() string {
	return o.Name()
}

func (o OS) IsCompatibleWith(os string) bool {
	for _, compatible := range o.CompatibleWith() {
		if os == compatible {
			return true
		}
	}
	return false
}

func (o *OS) Parse(name string) error {
	for _, os := range o.Supported() {
		if compare.ContainsLower(name, os.Name()) {
			*o = os
			return nil
		}
	}

	for _, os := range o.Supported() {
		for _, alias := range os.CompatibleWith() {
			if compare.ContainsLower(name, alias) {
				*o = os
				return nil
			}
		}
	}

	// Return an error if no match is found
	return fmt.Errorf("unable to parse OS from name: %s", name)
}

func (o *OS) Default() OS {
	return Linux
}
