package none

import (
	"github.com/hashicorp/go-getter/v2"

	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/sources/common"
	"github.com/idelchi/godyl/pkg/path/file"
)

// None represents a no-op source that performs no actual installation.
// It implements the Populator interface but all methods are no-ops.
type None struct{}

// Get returns an empty string for any attribute key.
func (n *None) Get(_ string) string {
	return ""
}

// Initialize is a no-op implementation of the Populator interface.
func (n *None) Initialize(_ string) error {
	return nil
}

// Version is a no-op implementation of the Populator interface.
func (n *None) Version(_ string) error {
	return nil
}

// Path is a no-op implementation of the Populator interface.
func (n *None) Path(_ string, _ []string, _ string, _ match.Requirements) error {
	return nil
}

// Install is a no-op implementation of the Populator interface.
// Returns empty values as no actual installation is performed.
func (n *None) Install(_ common.InstallData, _ getter.ProgressTracker) (string, file.File, error) {
	return "", file.File(""), nil
}
