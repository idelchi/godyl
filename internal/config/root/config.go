package root

import (
	"github.com/idelchi/godyl/internal/config/common"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/folder"
)

// Root holds the root level configuration options.
type Root struct {
	// Tokens store authentication tokens for various sources
	Tokens Tokens `json:"tokens" mapstructure:",squash"`

	// Inherit specifies the default scheme to inherit from when no scheme is specified
	Inherit string `json:"inherit" mapstructure:"inherit"`

	// ErrorFile specifies the file to log errors
	ErrorFile file.File `json:"error-file" mapstructure:"error-file"`

	// Tools specifies the tools file to be used
	Tools string `json:"tools" mapstructure:"tools"`

	// Defaults specifies the default file to be used
	Defaults file.File `json:"defaults" mapstructure:"defaults"`

	// ConfigFile specifies the configuration file to be used
	ConfigFile file.File `json:"config-file" mapstructure:"config-file"`

	// LogLevel specifies the logging level
	LogLevel string `json:"log-level" mapstructure:"log-level" validate:"oneof=silent debug info warn error always"`

	// EnvFile specifies the environment files to be used
	EnvFile []file.File `json:"env-file" mapstructure:"env-file"`

	// Cache holds the cache configuration options
	Cache Cache `json:"cache" mapstructure:",squash"`

	// Parallel specifies the number of parallel operations
	Parallel int `json:"parallel" mapstructure:"parallel" validate:"gte=0"`

	// Verbose specifies the verbosity level
	Verbose int `json:"verbose" mapstructure:"verbose"`

	// Show specifies the verbosity level for showing output
	Show Verbosity `json:"show" mapstructure:"show"`

	// NoVerifySSL disables SSL verification
	NoVerifySSL bool `json:"no-verify-ssl" mapstructure:"no-verify-ssl"`

	// NoProgress disables progress indicators
	NoProgress bool `json:"no-progress" mapstructure:"no-progress"`

	// Tracker embed the common tracker configuration, allowing to tracker
	// whether configuration values have been explicitly set or defaulted
	common.Tracker `json:"-" mapstructure:"-"`
}

// Cache holds the configuration options for caching.
type Cache struct {
	// Path to cache folder
	Dir folder.Folder `json:"cache-dir" mapstructure:"cache-dir"`

	// Disabled disables cache interaction
	Disabled bool `json:"no-cache" mapstructure:"no-cache"`
}

// Tokens holds the configuration options for authentication tokens.
type Tokens struct {
	// GitHub token for authentication
	GitHub string `json:"github-token" mapstructure:"github-token" mask:"fixed"`

	// GitLab token for authentication
	GitLab string `json:"gitlab-token" mapstructure:"gitlab-token" mask:"fixed"`

	// URL token for authentication
	URL string `json:"url-token" mapstructure:"url-token" mask:"fixed"`
}
