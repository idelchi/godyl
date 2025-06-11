package common

import (
	"github.com/idelchi/godyl/internal/debug"
	"github.com/idelchi/godyl/pkg/koanfx"
)

type Tracker struct {
	tracker *koanfx.Tracker `mapstructure:"-" yaml:"-"`
}

func (t *Tracker) Validate() error {
	return nil
}

func (t *Tracker) Store(tracker *koanfx.Tracker) {
	t.tracker = tracker
}

func (t *Tracker) Retrieve() *koanfx.Tracker {
	return t.tracker
}

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
	} else {
		debug.Debug("name %s found in tracker", name)
	}

	return t.tracker.IsSet(name)
}
