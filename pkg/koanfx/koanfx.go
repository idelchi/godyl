package koanfx

import (
	"fmt"

	"github.com/go-viper/mapstructure/v2"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/pflag"
)

// KoanfWithTracker is a wrapper around koanf.Koanf that tracks the keys that have been set.
type KoanfWithTracker struct {
	*koanf.Koanf
	flags   *pflag.FlagSet
	Tracker *Tracker

	active func()
}

// New creates a new KoanfWithTracker instance with a new koanf.Koanf instance.
func New() *KoanfWithTracker {
	return &KoanfWithTracker{
		Koanf: koanf.New("."),
	}
}

// NewWithTracker creates a new KoanfWithTracker instance with a new koanf.Koanf instance and a pflag.FlagSet.
func NewWithTracker(flags *pflag.FlagSet) *KoanfWithTracker {
	return &KoanfWithTracker{
		Koanf:   koanf.New("."),
		flags:   flags,
		Tracker: NewTracker(),
	}
}

// ResetFlags overwrites the flags in the KoanfWithTracker instance with a new pflag.FlagSet.
func (kwt *KoanfWithTracker) ResetFlags(newF *pflag.FlagSet) *KoanfWithTracker {
	return &KoanfWithTracker{
		Koanf:   kwt.Koanf,
		flags:   newF,
		Tracker: kwt.Tracker,
	}
}

// ResetTracker overwrites the Tracker in the KoanfWithTracker instance with a new Tracker.
func (kwt *KoanfWithTracker) ResetTracker() *KoanfWithTracker {
	return &KoanfWithTracker{
		Koanf:   kwt.Koanf,
		flags:   kwt.flags,
		Tracker: NewTracker(),
	}
}

// ResetKoanf overwrites the Koanf in the KoanfWithTracker instance with a new koanf.Koanf instance.
func (kwt *KoanfWithTracker) ResetKoanf(newK *koanf.Koanf) *KoanfWithTracker {
	if newK == nil {
		newK = koanf.New(".")
	}

	return &KoanfWithTracker{
		Koanf:   newK,
		flags:   kwt.flags,
		Tracker: kwt.Tracker,
	}
}

// IsSet checks if a key is set in the KoanfWithTracker instance.
func (kwt *KoanfWithTracker) IsSet(key string) bool {
	return kwt.Tracker.IsSet(key)
}

// TrackAll tracks all keys from the koanf instance.
func (kwt *KoanfWithTracker) TrackAll() *KoanfWithTracker {
	kwt.active = kwt.withAll()

	return kwt
}

// TrackFlags tracks changed flags from a FlagSet.
func (kwt *KoanfWithTracker) TrackFlags() *KoanfWithTracker {
	kwt.active = kwt.withFlags()

	return kwt
}

// Track sets the active tracking function to the last one set.
func (kwt *KoanfWithTracker) Track() *KoanfWithTracker {
	if kwt.active != nil {
		kwt.active()
	}

	return kwt
}

// Load loads configuration from a provider and tracks the source.
func (kwt *KoanfWithTracker) Load(p koanf.Provider, pa koanf.Parser, opts ...koanf.Option) error {
	if err := kwt.Koanf.Load(p, pa, opts...); err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	kwt.Track()

	return nil
}

// Unmarshal unmarshals the configuration into the provided struct.
func (kwt *KoanfWithTracker) Unmarshal(destination any, opts ...Option) error {
	return Unmarshal(kwt.Koanf, "", destination, opts...)
}

// UnmarshalWithMetadata unmarshals the configuration into the provided struct and returns metadata.
func (kwt *KoanfWithTracker) UnmarshalWithMetadata(destination any, opts ...Option) (mapstructure.Metadata, error) {
	return UnmarshalWithMetadata(kwt.Koanf, "", destination, opts...)
}

// AsMapAny unmarshals the configuration into a map[string]any.
func (kwt *KoanfWithTracker) AsMapAny() (map[string]any, error) {
	var mapAny map[string]any
	if err := kwt.Unmarshal(&mapAny); err != nil {
		return nil, fmt.Errorf("unmarshaling into struct: %w", err)
	}

	return mapAny, nil
}

// withAll tracks all keys from the koanf instance.
func (kwt *KoanfWithTracker) withAll() func() {
	return func() {
		kwt.Tracker.TrackAll(kwt.Koanf)
	}
}

// withFlags tracks changed flags from a FlagSet.
func (kwt *KoanfWithTracker) withFlags() func() {
	return func() {
		kwt.Tracker.TrackFlags(kwt.flags)
	}
}
