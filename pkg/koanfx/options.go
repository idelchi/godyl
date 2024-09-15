package koanfx

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/knadh/koanf/v2"
)

// NewDefaultUnmarshalConfig creates a new default UnmarshalConf with the default tag.
func NewDefaultUnmarshalConfig() koanf.UnmarshalConf {
	return koanf.UnmarshalConf{
		Tag: "mapstructure",
	}
}

// NewUnmarshalConfig creates a new UnmarshalConf with the default decoder config.
func NewUnmarshalConfig() koanf.UnmarshalConf {
	koanfig := NewDefaultUnmarshalConfig()

	koanfig.DecoderConfig = &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			textUnmarshalerHookFunc(),
		),
		WeaklyTypedInput: true,
		// IgnoreUntaggedFields: true,
	}

	return koanfig
}

// Option is a function that modifies an UnmarshalConf.
type Option func(*koanf.UnmarshalConf)

// WithErrorUnused causes errors on unused fields.
func WithErrorUnused() Option {
	return func(conf *koanf.UnmarshalConf) {
		if conf.DecoderConfig == nil {
			conf.DecoderConfig = NewUnmarshalConfig().DecoderConfig
		}

		conf.DecoderConfig.ErrorUnused = true
	}
}

// WithFlatPaths enables flat paths.
func WithFlatPaths() Option {
	return func(conf *koanf.UnmarshalConf) {
		conf.FlatPaths = true
	}
}

// WithSquash enables squash.
func WithSquash() Option {
	return func(conf *koanf.UnmarshalConf) {
		if conf.DecoderConfig == nil {
			conf.DecoderConfig = NewUnmarshalConfig().DecoderConfig
		}

		conf.DecoderConfig.Squash = true
	}
}
