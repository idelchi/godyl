package config

import (
	"errors"

	"github.com/spf13/viper"
)

// ErrUsage is returned when there is an error in the configuration.
var ErrUsage = errors.New("usage error")

// Embedded holds the embedded files for the application.
type Embedded struct {
	Defaults []byte
	Tools    []byte
	Template []byte
}

// Config holds all the configuration options for godyl.
type Config struct {
	// Root level configuration
	Root Root

	Tool Tool

	// Dump level configuration
	Dump Dump
}

// IsSet checks if a flag is set in viper.
func IsSet(flag string) bool {
	return viper.IsSet(flag)
}
