package main

import (
	"fmt"
	"slices"

	"github.com/go-playground/validator/v10"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/pkg/file"
	"github.com/idelchi/godyl/pkg/logger"
)

type Update struct {
	Strategy tools.Strategy `mapstructure:"strategy"`
	Update   bool           `mapstructure:"update"`
}

type Tokens struct {
	GitHub string `mapstructure:"github-token"`
	GitLab string `mapstructure:"gitlab-token"`
	URL    string `mapstructure:"url-token"`
}

type Show struct {
	Config   bool `mapstructure:"show-config"`
	Env      bool `mapstructure:"show-env"`
	Defaults bool `mapstructure:"show-defaults"`
}

// Config holds all the configuration options for godyl.
type Config struct {
	Tools    string
	Tags     []string
	Defaults file.File
	Update   Update `mapstructure:",squash"`
	Dry      bool
	Detect   bool
	Log      logger.Level
	DotEnv   file.File `mapstructure:"dot-env"`
	// Number of parallel downloads
	Parallel int `validate:"gte=0"`

	Show Show `mapstructure:",squash"`

	// Show help message
	Help bool
	// Show version information
	Version bool

	Output string
	Tokens Tokens `mapstructure:",squash"`
	Source sources.Type
}

func (c *Config) Validate() error {
	allowedUpdateStrategies := []tools.Strategy{tools.None, tools.Upgrade, tools.Force}
	if !slices.Contains(allowedUpdateStrategies, c.Update.Strategy) {
		return fmt.Errorf("invalid update strategy: %q: allowed are %v", c.Update.Strategy, allowedUpdateStrategies)
	}

	if IsSet("config") && !c.Defaults.Exists() {
		return fmt.Errorf("defaults file %q does not exist", c.Defaults)
	}

	if IsSet("dot-env") && !c.DotEnv.Exists() {
		return fmt.Errorf("dot-env file %q does not exist", c.DotEnv)
	}

	validate := validator.New()
	if err := validate.Struct(c); err != nil {
		return fmt.Errorf("validating config: %w", err)
	}
	return nil
}
