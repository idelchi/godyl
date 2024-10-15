package tools

import (
	"github.com/idelchi/godyl/internal/version"
)

type Strategy string

const (
	None    Strategy = "none"
	Upgrade Strategy = "upgrade"
	Force   Strategy = "force"
)

func (s Strategy) Check(t *Tool) error {
	if t.Strategy == None && t.Exists() {
		return ErrAlreadyExists
	}

	return nil
}

func (s Strategy) Upgrade(t *Tool) error {
	if !t.Exists() {
		return nil
	}

	switch t.Strategy {
	case None:
		return ErrAlreadyExists
	case Upgrade:
		exe := version.NewExecutable(t.Output, t.Exe.Name)
		err := exe.ParseVersion()
		if err != nil {
			return nil
		}
		if version, err := version.NewDefaultVersionParser().ParseString(t.Version); err == nil {
			if exe.Version == version {
				return ErrAlreadyExists
			}
		}

		return nil
	case Force:
		return nil
	default:
		return nil
	}
}
