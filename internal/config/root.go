package config

import (
	"fmt"

	"github.com/idelchi/godyl/internal/tools/sources/url"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/folder"
	"github.com/idelchi/godyl/pkg/validator"
)

// Root holds the root configuration options.
type Root struct {
	Tokens     Tokens    `mapstructure:",squash"`
	Log        string    `validate:"oneof=silent debug info warn error always"`
	ConfigFile file.File `mapstructure:"config-file"`
	Defaults   file.File
	EnvFile    []file.File   `mapstructure:"env-file"`
	Cache      CacheSettings `mapstructure:",squash"`
	Default    string

	viperable `json:"-" mapstructure:"-" yaml:"-"`
}

type CacheSettings struct {
	// Path to cache folder
	Dir folder.Folder `mapstructure:"cache-dir"`

	// Disabled disables cache interaction
	Disabled bool `mapstructure:"no-cache"`
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

	return validator.Validate(r)
}
