package koanfx

import (
	"github.com/knadh/koanf/v2"
	"github.com/spf13/pflag"
)

// SetMap is a map that tracks the keys that have been set and whether they were changed.
type SetMap map[string]bool

// Tracker is a structure that tracks the keys in a koanf instance and their changed status.
type Tracker struct {
	track SetMap
}

// NewTracker creates a new Tracker instance.
func NewTracker() *Tracker {
	return &Tracker{
		track: make(SetMap),
	}
}

// Names returns the names of all tracked keys.
func (dt *Tracker) Names() []string {
	names := make([]string, 0, len(dt.track))
	for key := range dt.track {
		names = append(names, key)
	}

	return names
}

// Exists checks if a key exists in the tracker.
func (dt *Tracker) Exists(key string) bool {
	_, exists := dt.track[key]

	return exists
}

// IsSet checks if a key is set in the tracker.
func (dt *Tracker) IsSet(key string) bool {
	return dt.track[key]
}

// set marks a key as set in the tracker.
func (dt *Tracker) set(key string) {
	dt.track[key] = true
}

// add adds a key to the tracker without marking it as set.
func (dt *Tracker) add(key string) {
	dt.track[key] = false
}

// TrackAll tracks all keys from the koanf instance.
func (dt *Tracker) TrackAll(k *koanf.Koanf) {
	for key := range k.All() {
		dt.set(key)
	}
}

// TrackFlags tracks changed flags from a FlagSet.
func (dt *Tracker) TrackFlags(flags *pflag.FlagSet) {
	flags.VisitAll(func(f *pflag.Flag) {
		if f.Changed {
			dt.set(f.Name)
		} else if !dt.Exists(f.Name) {
			dt.add(f.Name)
		}
	})
}
