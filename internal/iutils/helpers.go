package iutils

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/fatih/structs"

	"github.com/idelchi/godyl/internal/tools/tags"
	"github.com/idelchi/godyl/pkg/env"
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

// Map is a type alias for a map with string keys and any values.
type Map map[string]any

// StructToMap converts a struct to a Map, using the "json" tag for field names.
func StructToMap(s any) Map {
	str := structs.New(s)

	str.TagName = "json"

	return str.Map()
}

// Keys returns the keys of the Map as a slice of strings.
func (m Map) Keys() []string {
	return slices.Collect(maps.Keys(m))
}

// Values returns the values of the Map as a slice of any.
func (m Map) Values() []any {
	return slices.Collect(maps.Values(m))
}
