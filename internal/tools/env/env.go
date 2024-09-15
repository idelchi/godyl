package env

import (
	"github.com/goccy/go-yaml/ast"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/unmarshal"
)

// Values represents a map of string keys to any values.
type Env struct {
	env.Env `yaml:",inline"`
}

// UnmarshalYAML implements custom YAML unmarshaling for Exe configuration.
// Supports both scalar values (treated as executable name) and map values.
func (e *Env) UnmarshalYAML(node ast.Node) error {
	type raw Env

	return unmarshal.MapWithAppend(&e.Env, node)
}
