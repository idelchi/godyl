package koanfx

import (
	"fmt"

	"github.com/go-viper/mapstructure/v2"
	"github.com/knadh/koanf/v2"
)

// UnmarshalAll unmarshals all keys from the Koanf instance into the provided struct.
func UnmarshalAll(koanf *koanf.Koanf, destination any, opts ...Option) error {
	return Unmarshal(koanf, "", destination, opts...)
}

// UnmarshalAllWithMetadata unmarshals all keys from the Koanf instance into the provided struct and returns metadata.
func UnmarshalAllWithMetadata(koanf *koanf.Koanf, destination any, opts ...Option) (mapstructure.Metadata, error) {
	return UnmarshalWithMetadata(koanf, "", destination, opts...)
}

// Unmarshal unmarshals with optional additional options, starting from the provided key.
func Unmarshal(koanf *koanf.Koanf, key string, destination any, opts ...Option) error {
	conf := NewDefaultUnmarshalConfig()

	if len(opts) > 0 {
		// If options are provided, we need to create a new config
		conf = NewUnmarshalConfig()
		// And set the result to the provided value
		conf.DecoderConfig.Result = destination
	}

	// Apply all additional options
	for _, opt := range opts {
		opt(&conf)
	}

	if err := koanf.UnmarshalWithConf(key, destination, conf); err != nil {
		return fmt.Errorf("unmarshalling config: %w", err)
	}

	return nil
}

// UnmarshalWithMetadata unmarshals with metadata and optional additional options, starting from the provided key.
func UnmarshalWithMetadata(
	koanf *koanf.Koanf,
	key string,
	destination any,
	opts ...Option,
) (mapstructure.Metadata, error) {
	var metadata mapstructure.Metadata

	conf := NewUnmarshalConfig()

	conf.DecoderConfig.Metadata = &metadata
	conf.DecoderConfig.Result = destination

	// Apply all additional options
	for _, opt := range opts {
		opt(&conf)
	}

	if err := koanf.UnmarshalWithConf(key, destination, conf); err != nil {
		return metadata, fmt.Errorf("unmarshalling config: %w", err)
	}

	return metadata, nil
}
