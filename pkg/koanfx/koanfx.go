package koanfx

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/pflag"
)

type KoanfWithTracker struct {
	*koanf.Koanf
	flags   *pflag.FlagSet
	Tracker *Tracker

	active func()
}

func New() *KoanfWithTracker {
	return &KoanfWithTracker{
		Koanf: koanf.New("."),
	}
}

func NewWithTracker(flags *pflag.FlagSet) *KoanfWithTracker {
	return &KoanfWithTracker{
		Koanf:   koanf.New("."),
		flags:   flags,
		Tracker: NewTracker(),
	}
}

func (kwt *KoanfWithTracker) ResetFlags(flags *pflag.FlagSet) *KoanfWithTracker {
	return &KoanfWithTracker{
		Koanf:   kwt.Koanf,
		flags:   flags,
		Tracker: kwt.Tracker,
	}
}

func (kwt *KoanfWithTracker) ResetTracker() *KoanfWithTracker {
	return &KoanfWithTracker{
		Koanf:   kwt.Koanf,
		flags:   kwt.flags,
		Tracker: NewTracker(),
	}
}

func (kwt *KoanfWithTracker) ResetK(k *koanf.Koanf) *KoanfWithTracker {
	if k == nil {
		k = koanf.New(".")
	}

	return &KoanfWithTracker{
		Koanf:   k,
		flags:   kwt.flags,
		Tracker: kwt.Tracker,
	}
}

func (kwt *KoanfWithTracker) IsSet(key string) bool {
	return kwt.Tracker.IsSet(key)
}

func (kwt *KoanfWithTracker) withAll() func() {
	return func() {
		kwt.Tracker.TrackAll(kwt.Koanf)
	}
}

func (kwt *KoanfWithTracker) withFlags() func() {
	return func() {
		kwt.Tracker.TrackFlags(kwt.flags)
	}
}

func (kwt *KoanfWithTracker) TrackAll() *KoanfWithTracker {
	kwt.active = kwt.withAll()

	return kwt
}

func (kwt *KoanfWithTracker) TrackFlags() *KoanfWithTracker {
	kwt.active = kwt.withFlags()

	return kwt
}

func (kwt *KoanfWithTracker) Track() *KoanfWithTracker {
	if kwt.active != nil {
		kwt.active()
	}

	return kwt
}

// LoadWithTracking loads configuration from a provider and tracks the source.
func (kwt *KoanfWithTracker) Load(p koanf.Provider, pa koanf.Parser, opts ...koanf.Option) error {
	if err := kwt.Koanf.Load(p, pa, opts...); err != nil {
		return err
	}

	kwt.Track()

	return nil
}

func (kwt *KoanfWithTracker) Unmarshal(c any, opts ...Option) error {
	return Unmarshal(kwt.Koanf, "", c, opts...)
}

func (kwt *KoanfWithTracker) UnmarshalWithMetadata(v any, opts ...Option) (mapstructure.Metadata, error) {
	return UnmarshalWithMetadata(kwt.Koanf, "", v, opts...)
}
