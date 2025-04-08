package none

import (
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/sources/common"
	"github.com/idelchi/godyl/pkg/path/file"
)

type None struct{}

func (n *None) Get(_ string) string {
	return ""
}

func (n *None) Initialize(_ string) error {
	return nil
}

func (n *None) Exe() error {
	return nil
}

func (n *None) Version(_ string) error {
	return nil
}

func (n *None) Path(_ string, _ []string, _ string, _ match.Requirements) error {
	return nil
}

func (n *None) Install(_ common.InstallData) (string, file.File, error) {
	return "", file.File(""), nil
}
