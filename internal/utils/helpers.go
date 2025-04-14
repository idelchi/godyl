package utils

import (
	"fmt"
	"strings"

	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/path/file"

	"gopkg.in/yaml.v3"
)

// LoadFromDotEnv loads environment variables from a .env file.
func LoadDotEnv(path file.File) error {
	dotEnv, err := env.FromDotEnv(path.Expanded().Path())
	if err != nil {
		return fmt.Errorf("loading environment variables from %q: %w", path, err)
	}

	if dotEnv.Has("GODYL_CONFIG_FILE") {
		return fmt.Errorf("GODYL_CONFIG_FILE is not allowed in .env files")
	}

	if err := env.FromEnv().Normalized().Merged(dotEnv.Normalized()).ToEnv(); err != nil {
		return fmt.Errorf("setting environment variables: %w", err)
	}

	return nil
}

// Split splits tags into include and exclude lists.
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

// ParseBytes parses YAML bytes into a generic data structure.
func ParseBytes(yamlBytes []byte) any {
	var data any
	if err := yaml.Unmarshal(yamlBytes, &data); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	return data
}
