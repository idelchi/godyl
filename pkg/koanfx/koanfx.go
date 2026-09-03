package koanfx

import (
	"fmt"

	"github.com/go-viper/mapstructure/v2"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/pflag"
)

// Koanf wraps koanf.Koanf and records whether each loaded key came from explicit input or an untouched flag default.
type Koanf struct {
	// Koanf stores the merged configuration values.
	*koanf.Koanf

	// flags is the flag set inspected by TrackFlags.
	flags *pflag.FlagSet
	// Tracker stores explicit-input state for loaded keys.
	Tracker *Tracker

	// active is the tracking action applied after the next Load.
	active func()
}

// New creates a new Koanf instance with a new koanf.Koanf instance.
func New() *Koanf {
	return &Koanf{
		Koanf:   koanf.New("."),
		Tracker: NewTracker(),
	}
}

// FromStruct creates a new Koanf instance from a struct, loading the struct's fields into the koanf.Koanf
// instance.
func FromStruct(data any, tag string) (*Koanf, error) {
	k := New()
	if err := k.Load(structs.Provider(data, tag), nil); err != nil {
		return nil, err
	}

	return k, nil
}

// WithFlags overwrites the flags in the Koanf instance with a new pflag.FlagSet.
func (kwt *Koanf) WithFlags(newF *pflag.FlagSet) *Koanf {
	return &Koanf{
		Koanf:   kwt.Koanf,
		flags:   newF,
		Tracker: kwt.Tracker,
	}
}

// ClearTracker overwrites the Tracker in the Koanf instance with a fresh Tracker instance.
func (kwt *Koanf) ClearTracker() *Koanf {
	return &Koanf{
		Koanf:   kwt.Koanf,
		flags:   kwt.flags,
		Tracker: NewTracker(),
	}
}

// WithKoanf overwrites the Koanf in the Koanf instance with a new koanf.Koanf instance.
func (kwt *Koanf) WithKoanf(newK *koanf.Koanf) *Koanf {
	if newK == nil {
		newK = koanf.New(".")
	}

	return &Koanf{
		Koanf:   newK,
		flags:   kwt.flags,
		Tracker: kwt.Tracker,
	}
}

// IsSet checks if a key is set in the Koanf instance.
func (kwt *Koanf) IsSet(key string) bool {
	return kwt.Tracker.IsSet(key)
}

// TrackAll selects tracking of all keys produced by the next successful Load.
func (kwt *Koanf) TrackAll() *Koanf {
	kwt.active = kwt.withAll()

	return kwt
}

// TrackFlags selects tracking of changed flags and untouched defaults after the next successful Load.
func (kwt *Koanf) TrackFlags() *Koanf {
	kwt.active = kwt.withFlags()

	return kwt
}

// Track executes the tracking mode most recently selected by TrackAll or TrackFlags.
func (kwt *Koanf) Track() *Koanf {
	if kwt.active != nil {
		kwt.active()
	}

	return kwt
}

// Load merges a provider and applies the selected explicit-value tracking action.
func (kwt *Koanf) Load(p koanf.Provider, pa koanf.Parser, opts ...koanf.Option) error {
	if err := kwt.Koanf.Load(p, pa, opts...); err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	kwt.Track()

	return nil
}

// Unmarshal unmarshals the configuration into the provided struct.
func (kwt *Koanf) Unmarshal(destination any, opts ...Option) error {
	return Unmarshal(kwt.Koanf, "", destination, opts...)
}

// UnmarshalWithMetadata unmarshals the configuration into the provided struct and returns metadata.
func (kwt *Koanf) UnmarshalWithMetadata(destination any, opts ...Option) (mapstructure.Metadata, error) {
	return UnmarshalWithMetadata(kwt.Koanf, "", destination, opts...)
}

// Filtered returns a new Koanf instance that only contains the specified keys.
// Any keys that don't exist in the original instance are ignored.
func (kwt *Koanf) Filtered(keys ...string) *Koanf {
	// Create a new Koanf instance with the same delimiter
	filtered := New()

	// Iterate through the requested keys
	for _, key := range keys {
		// Only copy keys that exist in the original
		if kwt.Exists(key) {
			// Get the value and set it in the new instance
			if err := filtered.Set(key, kwt.Get(key)); err != nil {
				// Skip this key if setting fails
				continue
			}
		}
	}

	return filtered
}

// Map returns the koanf instance as a map.
func (kwt *Koanf) Map() map[string]any {
	return kwt.Raw()
}

// withAll tracks all keys from the koanf instance.
func (kwt *Koanf) withAll() func() {
	return func() {
		kwt.Tracker.TrackAll(kwt.Koanf)
	}
}

// withFlags tracks changed flags from a FlagSet.
func (kwt *Koanf) withFlags() func() {
	return func() {
		kwt.Tracker.TrackFlags(kwt.flags)
	}
}
