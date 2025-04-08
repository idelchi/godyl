package config

import (
	"fmt"

	"github.com/idelchi/godyl/internal/tools/sources/url"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/folder"
	"github.com/idelchi/godyl/pkg/validate"
)

// Root holds the root configuration options.
type Root struct {
	// Run without making any changes
	Dry bool

	// Log level (DEBUG, INFO, WARN, ERROR, SILENT)
	Log string `validate:"oneof=DEBUG INFO WARN ERROR SILENT"`

	// Path to config file
	ConfigFile file.File `mapstructure:"config-file"`

	// Path to .env file
	EnvFile []file.File `mapstructure:"env-file"`

	// Path to defaults file
	Defaults file.File

	// Cache settings
	Cache CacheSettings `mapstructure:",squash"`

	// Tokens for authentication
	Tokens Tokens `mapstructure:",squash"`

	// Viper instance
	viperable `mapstructure:"-" yaml:"-" json:"-"`
}

type CacheSettings struct {
	// Path to cache folder
	Dir folder.Folder `mapstructure:"cache-dir"`
	// Type of cache (file, sqlite)
	Type string `mapstructure:"cache-type" validate:"oneof=file sqlite"`
}

// Tokens holds the configuration options for authentication tokens.
type Tokens struct {
	// GitHub token for authentication
	GitHub string `mapstructure:"github-token" mask:"fixed"`
	// GitLab token for authentication
	GitLab string `mapstructure:"gitlab-token" mask:"fixed"`
	// URL token for authentication
	URL url.Token `mapstructure:",squash"`
}

// Validate checks the configuration for errors.
func (r *Root) Validate() error {
	if r.IsSet("config-file") && !r.ConfigFile.Expanded().Exists() {
		return fmt.Errorf("%w: config file %q does not exist", ErrUsage, r.ConfigFile.Expanded())
	}

	if r.IsSet("defaults") && !r.Defaults.Expanded().Exists() {
		return fmt.Errorf("%w: defaults file %q does not exist", ErrUsage, r.Defaults.Expanded())
	}

	if r.IsSet("env-file") {
		for _, file := range r.EnvFile {
			if !file.Expanded().Exists() {
				return fmt.Errorf("%w: env-file file %q does not exist", ErrUsage, file.Expanded())
			}
		}
	}

	return validate.Validate(r)
}
