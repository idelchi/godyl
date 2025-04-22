// Package strategy provides functionality for managing tool installation strategies.
package strategy

import (
	"fmt"

	"github.com/idelchi/godyl/internal/tools/result"
	"github.com/idelchi/godyl/pkg/version"
)

// Strategy represents the strategy for handling tool installation.
type Strategy string

const (
	// None indicates no strategy, meaning no action will be taken if the tool already exists.
	None Strategy = "none"
	// Sync indicates that the tool should only be modified if different from the desired state.
	Sync = "sync"
	// Force indicates that the tool should be installed or updated regardless of its current state.
	Force = "force"
)

func (s Strategy) String() string {
	return string(s)
}

// Tool represents the interface required by the Strategy type.
type Tool interface {
	Exists() bool
	GetCurrentVersion() string
	GetStrategy() Strategy
	GetTargetVersion() string
}

// Sync checks if the tool should be synced based on the strategy and its current version.
// It compares the existing version with the desired version and returns an error if the tool is already up to date.
func (s Strategy) Sync(t Tool) result.Result {
	// If the tool does not exist, no sync checks are necessary.
	if !t.Exists() {
		return result.WithOK("tool does not exist")
	}

	currentVersion := t.GetCurrentVersion()

	switch t.GetStrategy() {
	case None:
		// If the strategy is "None" and the tool exists, return an error indicating it already exists.
		return result.WithSkipped("already exists")
	case Sync:
		if currentVersion == "" {
			return result.WithOK("current version not retrievable, forcing update")
		}

		targetVersion := t.GetTargetVersion()

		// Convert versions for comparison
		source := version.To(currentVersion)
		target := version.To(targetVersion)

		// Check for conversion failures
		if source == nil {
			return result.WithOK(fmt.Sprintf("converting source version %q: failed: %q -> %q",
				currentVersion, currentVersion, targetVersion))
		}

		if target == nil {
			return result.WithOK(fmt.Sprintf(
				"converting target version %q: failed: %q -> %q",
				targetVersion,
				currentVersion,
				targetVersion,
			))
		}

		// Compare versions and return appropriate error
		if source.Equal(target) {
			return result.WithSkipped(fmt.Sprintf("current version %q and target version %q match", currentVersion, targetVersion))
		}

		return result.WithOK(fmt.Sprintf("current version %q and target version %q do not match", currentVersion, targetVersion))
	case Force:
		// If the strategy is "Force", always proceed with the installation or update.
		return result.WithOK("strategy is force, proceeding with sync")
	default:
		return result.WithFailed(fmt.Sprintf("unknown strategy %q", t.GetStrategy()))
	}
}
