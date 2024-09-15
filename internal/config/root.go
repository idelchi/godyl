package config

import (
	"fmt"

	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/folder"
)

// Root holds the root configuration options.
type Root struct {
	Tokens     Tokens        `mapstructure:",squash"`
	LogLevel   string        `mapstructure:"log-level" validate:"oneof=silent debug info warn error always"`
	ConfigFile file.File     `mapstructure:"config-file" yaml:"config-file"`
	Defaults   file.File     `mapstructure:"defaults"`
	EnvFile    []file.File   `mapstructure:"env-file" yaml:"env-file"`
	Cache      CacheSettings `mapstructure:",squash"`
	Inherit    string        `mapstructure:"inherit"`

	Parallel    int  `mapstructure:"parallel" validate:"gte=0"`
	NoVerifySSL bool `mapstructure:"no-verify-ssl" yaml:"no-verify-ssl"`
	NoProgress  bool `mapstructure:"no-progress" yaml:"no-progress"`

	ErrorFile file.File `mapstructure:"error-file" yaml:"error-file"`
	Detailed  bool

	Show bool `mapstructure:"show"`

	trackable `json:"-" mapstructure:"-" yaml:"-"`
}

func (r *Root) Validate() error {
	if r.IsSet("defaults") && !r.Defaults.Exists() {
		return fmt.Errorf("%w: %q does not exist", ErrUsage, r.Defaults)
	}

	return nil
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
	URL Token `mapstructure:",squash"`
}

// Token contains authentication configuration for URL requests.
type Token struct {
	// Token is the authentication token value.
	Token string `mapstructure:"url-token" mask:"fixed"`

	// Header is the HTTP header name for the token.
	Header string `mapstructure:"url-token-header"`

	// Scheme is the authentication scheme (e.g., "Bearer").
	Scheme string `mapstructure:"url-token-scheme"`
}
