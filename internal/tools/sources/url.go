package sources

import (
	"github.com/idelchi/godyl/internal/match"
)

type URL struct {
	URL   string
	Token string

	Data Metadata `yaml:"-"`
}

func (u *URL) Get(attribute string) string {
	return u.Data.Get(attribute)
}

func (u *URL) Initialize(name string) error {
	return nil
}

func (u *URL) Exe() error {
	return nil
}

func (u *URL) Version(name string) error {
	return nil
}

func (u *URL) Path(name string, _ []string, _ string, _ match.Requirements) error {
	u.Data.Set("path", name)

	return nil
}

func (u *URL) Install(d InstallData) (output string, err error) {
	return Download(d)
}
