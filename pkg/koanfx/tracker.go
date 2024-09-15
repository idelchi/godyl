package koanfx

import (
	"github.com/knadh/koanf/v2"
	"github.com/spf13/pflag"
)

// SetMap is a map that tracks the keys that have been set and whether they were changed.
type SetMap map[string]bool

type Tracker struct {
	track SetMap
}

func NewTracker() *Tracker {
	return &Tracker{
		track: make(SetMap),
	}
}

func (dt *Tracker) Names() []string {
	names := make([]string, 0, len(dt.track))
	for key := range dt.track {
		names = append(names, key)
	}

	return names
}

func (dt *Tracker) Exists(key string) bool {
	_, exists := dt.track[key]
	return exists
}

func (dt *Tracker) IsSet(key string) bool {
	return dt.track[key]
}

func (dt *Tracker) set(key string) {
	dt.track[key] = true
}

func (dt *Tracker) add(key string) {
	dt.track[key] = false
}

// TrackAll tracks all keys from the koanf instance
func (dt *Tracker) TrackAll(k *koanf.Koanf) {
	for key := range k.All() {
		dt.set(key)
	}
}

// TrackFlags tracks changed flags from a FlagSet
func (dt *Tracker) TrackFlags(flags *pflag.FlagSet) {
	flags.VisitAll(func(f *pflag.Flag) {
		if f.Changed {
			dt.set(f.Name)
		} else if !dt.Exists(f.Name) {
			dt.add(f.Name)
		}
	})
}
