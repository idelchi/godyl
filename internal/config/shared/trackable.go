package shared

import (
	"github.com/idelchi/godyl/internal/debug"
	"github.com/idelchi/godyl/pkg/koanfx"
)

// Tracker stores whether configuration values were supplied explicitly rather than inherited from defaults.
type Tracker struct {
	// tracker holds the explicit-value state collected while Koanf layers are loaded.
	tracker *koanfx.Tracker `mapstructure:"-" yaml:"-"`
}

// Validate satisfies the configuration validation contract; tracker state has no independently invalid form.
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

// IsSet reports whether name was supplied explicitly and panics before Store initializes the tracker.
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
