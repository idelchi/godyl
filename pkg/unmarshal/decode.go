// Package unmarshal provides utilities for unmarshalling YAML data that can
// represent either a single item or a slice of items using goccy/go-yaml.
// It includes a generic type `SingleOrSlice` to handle this pattern, allowing
// flexible unmarshalling from YAML input.
//
// The package also provides functions to decode YAML nodes while optionally
// enforcing that only known fields are present, improving error handling
// in case of unexpected fields in the YAML input.
package unmarshal

import (
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

func Decode(node ast.Node, out any, allowUnknownFields ...bool) error {
	opt := yaml.Strict()

	if len(allowUnknownFields) > 0 && allowUnknownFields[0] {
		// Disable strict mode to allow unknown fields
		opt = nil
	}

	if err := yaml.NodeToValue(node, out, opt); err != nil {
		return err
	}

	return nil
}
