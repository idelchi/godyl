package processor

import (
	"github.com/hashicorp/go-getter/v2"

	"github.com/idelchi/godyl/pkg/download"
	"github.com/idelchi/godyl/pkg/download/progress"
)

// progressMgr wraps progress tracking functionality.
type progressMgr struct {
	trackable progress.Trackable
}

// newProgressMgr creates a new progressMgr.
func newProgressMgr(noProgress bool) *progressMgr {
	var trackable progress.Trackable

	if noProgress {
		trackable = progress.NewNoop()
	} else {
		trackable = progress.New(download.DefaultTimeout)
	}

	return &progressMgr{
		trackable: trackable,
	}
}

// Start starts the progress tracking.
func (m *progressMgr) Start() {
	m.trackable.Start()
}

// Wait waits for all progress tracking to complete.
func (m *progressMgr) Wait() {
	m.trackable.Wait()
}

// Tracker returns the underlying progress tracker.
func (m *progressMgr) Tracker() getter.ProgressTracker {
	return m.trackable
}
