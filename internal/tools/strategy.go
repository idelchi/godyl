package tools

import (
	"fmt"

	"github.com/idelchi/godyl/internal/executable"
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
	switch t.Strategy {
	case None:
		if t.Exists() {
			return ErrAlreadyExists
		}
	case Upgrade:
		if t.Exists() {
			exe := executable.New(t.Output, t.Exe)
			err := exe.ParseVersion()
			if err != nil {
				fmt.Printf("parsing version: %v\n", err)
			}
			if version, err := executable.NewDefaultVersionParser().ParseString(t.Version); err == nil {
				if exe.Version == version {
					return ErrAlreadyExists
				}
			}
		}

		return nil
	case Force:
		return nil
	default:
		return nil
	}

	return nil
}
