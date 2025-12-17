//nolint:revive // Package name is appropriate
package shared

import (
	"github.com/idelchi/godyl/internal/debug"
	"github.com/idelchi/godyl/pkg/koanfx"
)

// Tracker manages tracking of configuration values and their sources.
type Tracker struct {
	tracker *koanfx.Tracker `mapstructure:"-" yaml:"-"`
}

// Validate performs validation on the tracker configuration.
func (t *Tracker) Validate() error {
	return nil
}

// Store saves the provided koanfx.Tracker instance.
func (t *Tracker) Store(tracker *koanfx.Tracker) {
	t.tracker = tracker
}

// Retrieve returns the stored koanfx.Tracker instance.
func (t *Tracker) Retrieve() *koanfx.Tracker {
	return t.tracker
}

// IsSet checks if a configuration value with the given name has been explicitly set.
func (t *Tracker) IsSet(name string) bool {
	if t.tracker == nil {
		panic("tracker is not initialized")
	}

	if !t.tracker.Exists(name) {
		debug.Debug("name %s not found in tracker", name)
		debug.Debug("available names:")

		for _, n := range t.tracker.Names() {
			debug.Debug("- %s", n)
		}
	}

	return t.tracker.IsSet(name)
}
