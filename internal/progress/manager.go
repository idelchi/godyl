// Package progress provides progress tracking functionality.
package progress

import (
	"github.com/hashicorp/go-getter/v2"

	"github.com/idelchi/godyl/pkg/download/progress"
)

// Manager defines the interface for progress tracking operations.
type Manager interface {
	// Start starts the progress tracking.
	Start()
	// Wait waits for all progress tracking to complete.
	Wait()
	// Tracker returns the underlying progress tracker.
	Tracker() getter.ProgressTracker
}

// DefaultManager wraps the existing progress.Trackable.
type DefaultManager struct {
	trackable progress.Trackable
}

// NewDefaultManager creates a new DefaultManager.
func NewDefaultManager(noProgress bool) *DefaultManager {
	var trackable progress.Trackable

	if noProgress {
		trackable = progress.NewNoop()
	} else {
		trackable = progress.New()
	}

	return &DefaultManager{
		trackable: trackable,
	}
}

// Start starts the progress tracking.
func (m *DefaultManager) Start() {
	m.trackable.Start()
}

// Wait waits for all progress tracking to complete.
func (m *DefaultManager) Wait() {
	m.trackable.Wait()
}

// Tracker returns the underlying progress tracker.
//
//nolint:ireturn // Returns interface for flexibility - allows different progress tracker implementations
func (m *DefaultManager) Tracker() getter.ProgressTracker {
	return m.trackable
}
