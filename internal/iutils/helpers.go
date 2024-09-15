package iutils

import (
	"fmt"
	"strings"

	"github.com/idelchi/godyl/internal/tools/tags"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/koanfx"
	"github.com/idelchi/godyl/pkg/path/file"
)

// LoadDotEnv loads environment variables from a .env file.
func LoadDotEnv(path file.File) (env.Env, error) {
	dotEnv, err := env.FromDotEnv(path.Expanded().Path())
	if err != nil {
		return nil, fmt.Errorf("loading environment variables from %q: %w", path, err)
	}

	return dotEnv, nil
}

// SplitTags splits tags into include and exclude lists.
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

// Any returns the first non-zero value from the provided arguments.
func Any[T comparable](args ...T) T {
	var zero T

	for _, arg := range args {
		if arg != zero {
			return arg
		}
	}

	return zero
}

// StructToKoanf converts a struct to a `koanf` instance, using the "yaml" tag for field names.
func StructToKoanf(s any) (*koanfx.Koanf, error) {
	return koanfx.FromStruct(s, "yaml")
}
