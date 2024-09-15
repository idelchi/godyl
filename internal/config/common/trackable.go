package common

import (
	"github.com/idelchi/godyl/internal/debug"
	"github.com/idelchi/godyl/pkg/koanfx"
)

type Trackable struct {
	tracker *koanfx.Tracker `json:"-" mapstructure:"-"`
}

func (t *Trackable) Validate() error {
	return nil
}

func (t *Trackable) Store(tracker *koanfx.Tracker) {
	t.tracker = tracker
}

func (t *Trackable) Retrieve() *koanfx.Tracker {
	return t.tracker
}

func (t *Trackable) IsSet(name string) bool {
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
