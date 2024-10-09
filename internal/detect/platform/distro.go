package platform

import (
	"fmt"

	"github.com/idelchi/godyl/pkg/compare"
)

// Distro enum
type Distribution string

// Distro enum values
const (
	Debian  Distribution = "debian"
	Ubuntu  Distribution = "ubuntu"
	CentOS  Distribution = "centos"
	RedHat  Distribution = "redhat"
	Arch    Distribution = "arch"
	Alpine  Distribution = "alpine"
	Rasbian Distribution = "rasbian"
)

func (d Distribution) Supported() []Distribution {
	return []Distribution{Debian, Ubuntu, CentOS, RedHat, Arch, Alpine, Rasbian}
}

func (d *Distribution) From(distribution string) error {
	for _, distro := range d.Supported() {
		if compare.Lower(distro.String(), distribution) {
			*d = distro

			return nil
		}
	}

	return fmt.Errorf("%w: distribution %q", ErrNotFound, distribution)
}

func (d Distribution) String() string {
	return string(d)
}

func (d Distribution) Default() Distribution {
	return Distribution("")
}
