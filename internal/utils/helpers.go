package utils

import (
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/idelchi/godyl/internal/tools/tags"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/path/file"
)

// LoadFromDotEnv loads environment variables from a .env file.
func LoadDotEnv(path file.File) (env.Env, error) {
	dotEnv, err := env.FromDotEnv(path.Expanded().Path())
	if err != nil {
		return nil, fmt.Errorf("loading environment variables from %q: %w", path, err)
	}

	return dotEnv.Normalized(), nil
}

// Split splits tags into include and exclude lists.
func SplitTags(tagList []string) tags.IncludeTags {
	tags := tags.IncludeTags{}

	for _, tag := range tagList {
		if strings.HasPrefix(tag, "!") {
			tags.Exclude = append(tags.Exclude, tag[1:])
		} else {
			tags.Include = append(tags.Include, tag)
		}
	}

	return tags
}

// ParseBytes parses YAML bytes into a generic data structure.
func ParseBytes(yamlBytes []byte) any {
	var data any
	if err := yaml.Unmarshal(yamlBytes, &data); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	return data
}
