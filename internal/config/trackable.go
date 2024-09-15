package config

import (
	"github.com/idelchi/godyl/internal/debug"
	"github.com/idelchi/godyl/pkg/koanfx"
)

type trackable struct {
	tracker *koanfx.Tracker `json:"-" mapstructure:"-" yaml:"-"`
}

func (t *trackable) Validate() error {
	return nil
}

func (t *trackable) StoreTracker(tracker *koanfx.Tracker) {
	t.tracker = tracker
}

func (t *trackable) RetrieveTracker() *koanfx.Tracker {
	return t.tracker
}

func (t *trackable) IsSet(name string) bool {
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
