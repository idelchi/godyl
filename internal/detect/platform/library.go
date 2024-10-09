package platform

import (
	"fmt"

	"github.com/idelchi/godyl/pkg/compare"
)

type Library string

const (
	GNU        Library = "gnu"
	Musl       Library = "musl"
	MSVC       Library = "msvc"
	LibAndroid Library = "android"
)

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
		return GNU
	}
}

func (l Library) Supported() []Library {
	return []Library{Musl, GNU, MSVC, LibAndroid}
}

func (l *Library) From(library string) error {
	for _, lib := range l.Supported() {
		if compare.Lower(library, lib.Name()) {
			*l = lib

			return nil
		}
	}

	for _, lib := range l.Supported() {
		if lib.IsCompatibleWith(library) {
			*l = lib

			return nil
		}
	}

	*l = ""

	return nil
}

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

func (l Library) IsCompatibleWith(lib string) bool {
	for _, compatible := range l.CompatibleWith() {
		if lib == compatible {
			return true
		}
	}
	return false
}

func (l Library) Name() string {
	return string(l)
}

func (l Library) String() string {
	return l.Name()
}

func (l *Library) Parse(name string) error {
	for _, library := range l.Supported() {
		if compare.ContainsLower(name, library.Name()) {
			*l = library
			return nil
		}
	}

	for _, library := range l.Supported() {
		for _, alias := range library.CompatibleWith() {
			if compare.ContainsLower(name, alias) {
				*l = library
				return nil
			}
		}
	}

	// Return an error if no match is found
	return fmt.Errorf("unable to parse library from name: %s", name)
}
