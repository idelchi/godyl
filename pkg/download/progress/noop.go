package progress

import "io"

// Noop is a dummy progress tracker that implements the Trackable interface
// but performs no tracking. Useful when progress display is disabled.
type Noop struct{}

// NewNoop returns a new instance of Noop.
func NewNoop() *Noop {
	return &Noop{}
}

// Start does nothing.
func (n *Noop) Start() {}

// Wait does nothing.
func (n *Noop) Wait() {}

// TrackProgress returns the unmodified stream without wrapping.
func (n *Noop) TrackProgress(_ string, _, _ int64, rc io.ReadCloser) io.ReadCloser {
	return rc
}
