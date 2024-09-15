package koanfx

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/knadh/koanf/v2"
)

func UnmarshalAll(k *koanf.Koanf, c any, opts ...Option) error {
	return Unmarshal(k, "", c, opts...)
}

func UnmarshalAllWithMetadata(k *koanf.Koanf, v any, opts ...Option) (mapstructure.Metadata, error) {
	return UnmarshalWithMetadata(k, "", v, opts...)
}

// Unmarshal unmarshals with optional additional options
func Unmarshal(k *koanf.Koanf, key string, c any, opts ...Option) error {
	conf := NewDefaultUnmarshalConfig()

	if len(opts) > 0 {
		// If options are provided, we need to create a new config
		conf = NewUnmarshalConfig()
		// And set the result to the provided value
		conf.DecoderConfig.Result = c
	}

	// Apply all additional options
	for _, opt := range opts {
		opt(&conf)
	}

	return k.UnmarshalWithConf(key, c, conf)
}

// UnmarshalWithMetadata unmarshals with metadata and optional additional options
func UnmarshalWithMetadata(k *koanf.Koanf, key string, c any, opts ...Option) (mapstructure.Metadata, error) {
	var md mapstructure.Metadata

	conf := NewUnmarshalConfig()

	conf.DecoderConfig.Metadata = &md
	conf.DecoderConfig.Result = c

	// Apply all additional options
	for _, opt := range opts {
		opt(&conf)
	}

	return md, k.UnmarshalWithConf(key, c, conf)
}
