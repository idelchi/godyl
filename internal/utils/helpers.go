package utils

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/file"
	"github.com/idelchi/godyl/pkg/logger"
)

// LoadDotEnv loads environment variables from a .env file.
func LoadDotEnv(path file.File) error {
	dotEnv, err := env.FromDotEnv(path.Name())
	if err != nil {
		return fmt.Errorf("loading environment variables from %q: %w", path.Name(), err)
	}

	env := env.FromEnv().Normalized().Merged(dotEnv.Normalized())

	if err := env.ToEnv(); err != nil {
		return fmt.Errorf("setting environment variables: %w", err)
	}

	return nil
}

// LoadTools loads the tools configuration.
func LoadTools(path string, log *logger.Logger) (tools.Tools, error) {
	var toolsList tools.Tools

	if err := toolsList.Load(path); err != nil {
		return toolsList, fmt.Errorf("loading tools from %q: %w", path, err)
	}

	log.Info("loaded %d tools from %q", len(toolsList), path)

	return toolsList, nil
}

// SplitTags splits tags into include and exclude lists.
func SplitTags(tags []string) ([]string, []string) {
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

// PrintYAMLBytes parses YAML bytes into a generic data structure.
func PrintYAMLBytes(yamlBytes []byte) any {
	var data any
	if err := yaml.Unmarshal(yamlBytes, &data); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	return data
}

// ValidateInput validates the command-line arguments.
func ValidateInput(toolsPath *string, args []string) error {
	switch len(args) {
	case 0:
		*toolsPath = "tools.yml"
	case 1:
		*toolsPath = args[0]
	}

	return nil
}

// FileExists checks if a file exists at the given path.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
