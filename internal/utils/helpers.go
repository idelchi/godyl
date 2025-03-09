package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/file"
	"github.com/idelchi/godyl/pkg/logger"

	"gopkg.in/yaml.v3"
)

// EnvironmentLoader loads environment variables from various sources.
type EnvironmentLoader interface {
	LoadFromDotEnv(path string) error
}

// ToolsLoader loads tools configuration from a file.
type ToolsLoader interface {
	LoadTools(path string, log *logger.Logger) (tools.Tools, error)
}

// DefaultEnvironmentLoader is the default implementation of EnvironmentLoader.
type DefaultEnvironmentLoader struct{}

// DefaultToolsLoader is the default implementation of ToolsLoader.
type DefaultToolsLoader struct{}

// LoadDotEnv loads environment variables from a .env file.
func LoadDotEnv(path file.File) error {
	loader := &DefaultEnvironmentLoader{}

	return loader.LoadFromDotEnv(path.Name())
}

// LoadFromDotEnv loads environment variables from a .env file.
func (l *DefaultEnvironmentLoader) LoadFromDotEnv(path string) error {
	dotEnv, err := env.FromDotEnv(path)
	if err != nil {
		return fmt.Errorf("loading environment variables from %q: %w", path, err)
	}

	environment := env.FromEnv().Normalized().Merged(dotEnv.Normalized())

	if err := environment.ToEnv(); err != nil {
		return fmt.Errorf("setting environment variables: %w", err)
	}

	return nil
}

// LoadTools loads the tools configuration.
func LoadTools(path file.File, log *logger.Logger) (tools.Tools, error) {
	loader := &DefaultToolsLoader{}

	return loader.LoadTools(string(path), log)
}

// LoadTools loads the tools configuration.
func (l *DefaultToolsLoader) LoadTools(path string, log *logger.Logger) (tools.Tools, error) {
	var toolsList tools.Tools

	if err := toolsList.Load(path); err != nil {
		return toolsList, fmt.Errorf("loading tools from %q: %w", path, err)
	}

	log.Info("loaded %d tools from %q", len(toolsList), path)

	return toolsList, nil
}

// TagSplitter splits tags into include and exclude lists.
type TagSplitter interface {
	Split(tags []string) (include, exclude []string)
}

// DefaultTagSplitter is the default implementation of TagSplitter.
type DefaultTagSplitter struct{}

// SplitTags splits tags into include and exclude lists.
func SplitTags(tags []string) ([]string, []string) {
	splitter := &DefaultTagSplitter{}

	return splitter.Split(tags)
}

// Split splits tags into include and exclude lists.
func (s *DefaultTagSplitter) Split(tags []string) ([]string, []string) {
	var withTags, withoutTags []string

	for _, tag := range tags {
		if strings.HasPrefix(tag, "!") {
			withoutTags = append(withoutTags, tag[1:])
		} else {
			withTags = append(withTags, tag)
		}
	}

	return withTags, withoutTags
}

// YAMLParser parses YAML data.
type YAMLParser interface {
	ParseBytes(yamlBytes []byte) any
}

// DefaultYAMLParser is the default implementation of YAMLParser.
type DefaultYAMLParser struct{}

// PrintYAMLBytes parses YAML bytes into a generic data structure.
func PrintYAMLBytes(yamlBytes []byte) any {
	parser := &DefaultYAMLParser{}

	return parser.ParseBytes(yamlBytes)
}

// ParseBytes parses YAML bytes into a generic data structure.
func (p *DefaultYAMLParser) ParseBytes(yamlBytes []byte) any {
	var data any
	if err := yaml.Unmarshal(yamlBytes, &data); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	return data
}

// InputValidator validates command-line arguments.
type InputValidator interface {
	ValidateInput(toolsPath *string, args []string) error
}

// DefaultInputValidator is the default implementation of InputValidator.
type DefaultInputValidator struct{}

// ValidateInput validates the command-line arguments.
func ValidateInput(toolsPath *string, args []string) error {
	validator := &DefaultInputValidator{}

	return validator.ValidateInput(toolsPath, args)
}

// ValidateInput validates the command-line arguments.
func (v *DefaultInputValidator) ValidateInput(toolsPath *string, args []string) error {
	switch len(args) {
	case 0:
		*toolsPath = "tools.yml"
	case 1:
		*toolsPath = args[0]
	}

	return nil
}

// FileChecker checks if a file exists.
type FileChecker interface {
	Exists(path string) bool
}

// DefaultFileChecker is the default implementation of FileChecker.
type DefaultFileChecker struct{}

// FileExists checks if a file exists at the given path.
func FileExists(path string) bool {
	checker := &DefaultFileChecker{}

	return checker.Exists(path)
}

// Exists checks if a file exists at the given path.
func (c *DefaultFileChecker) Exists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}
