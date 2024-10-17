package platform

import (
	"fmt"

	"github.com/idelchi/godyl/pkg/utils"
)

// Distribution represents a Linux distribution, with several predefined values.
type Distribution string

// Predefined Distribution values.
const (
	Debian  Distribution = "debian"  // Debian distribution
	Ubuntu  Distribution = "ubuntu"  // Ubuntu distribution
	CentOS  Distribution = "centos"  // CentOS distribution
	RedHat  Distribution = "redhat"  // RedHat distribution
	Arch    Distribution = "arch"    // Arch Linux distribution
	Alpine  Distribution = "alpine"  // Alpine Linux distribution
	Rasbian Distribution = "rasbian" // Rasbian distribution (used on Raspberry Pi)
)

// Available returns a slice of all available distributions.
func (d Distribution) Available() []Distribution {
	return []Distribution{Debian, Ubuntu, CentOS, RedHat, Arch, Alpine, Rasbian}
}

// From sets the Distribution based on the provided string, if it matches any available distribution.
func (d *Distribution) From(distribution string) error {
	for _, distro := range d.Available() {
		if utils.EqualLower(distro.String(), distribution) {
			*d = distro
			return nil
		}
	}

	return fmt.Errorf("%w: distribution %q", ErrNotFound, distribution)
}

// String returns the Distribution as a string.
func (d Distribution) String() string {
	return string(d)
}

// Default returns the default Distribution, which is an empty string.
func (d Distribution) Default() Distribution {
	return Distribution("")
}
