package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/tools/sources"
)

// Config holds all the configuration options for godyl.
type Config struct {
	// Defaults for tools. Allows setting a default subset of values for tools
	Defaults tools.Defaults

	// Path to file to load tools from
	Tools string

	// Tags to consider when selecting tools
	Tags []string

	// Config file to load
	Config string

	// Show help message
	Help bool
	// Show parsed configuration
	Show bool
	// Show version information
	Version bool

	// Number of parallel downloads
	Parallel int `validate:"gte=0"`
}

// Default method sets the default configuration values for godyl.
func (c *Config) Default() {
	dirname, err := os.UserHomeDir()
	if err != nil {
		dirname = "/tmp"
	}

	c.Defaults = tools.Defaults{
		Output: filepath.Join(dirname, ".local", "bin"),
		Source: sources.Source{
			Type: "github",
		},
		Hints: []match.Hint{
			{
				Pattern: `{{ .Exe }}`,
				Weight:  1,
			},
			// {
			// 	Pattern: "static",
			// 	Weight:  1,
			// },
			// {
			// 	Pattern: `.*{{ if eq .Platform.OS "linux" }}\.tar\.gz{{ else }}\.zip{{ end }}$`,
			// 	Weight:  1,
			// 	Regex:   true,
			// },
			// {
			// 	Pattern:        `musl`,
			// 	WeightTemplate: `{{ if and (eq .Platform.OS "linux") (eq .Platform.Distribution "alpine") }}1{{ else }}0{{ end }}`,
			// },
			// {
			// 	Pattern:        `gnu`,
			// 	WeightTemplate: `{{ if and (eq .Platform.OS "linux") (ne .Platform.Distribution "alpine") }}1{{ else }}0{{ end }}`,
			// },
		},
	}
}

// Validate the configuration.
func (c *Config) Validate() error {
	validate := validator.New()
	if err := validate.Struct(c); err != nil {
		return fmt.Errorf("validating config: %w", err)
	}
	return nil
}
