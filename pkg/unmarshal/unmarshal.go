package unmarshal

import (
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

// Decode decodes a YAML node into the specified output type, disallowing unknown fields.
func Decode(node ast.Node, out any) error {
	if err := yaml.NodeToValue(node, out, yaml.Strict()); err != nil {
		return err
	}

	return nil
}

// Strict unmarshals YAML data into the specified output type, disallowing unknown fields.
func Strict(data []byte, out any, options ...yaml.DecodeOption) error {
	options = append(options, yaml.Strict())

	if err := yaml.UnmarshalWithOptions(data, out, options...); err != nil {
		return err
	}

	return nil
}

// Lax unmarshals YAML data into the specified output type, allowing unknown fields.
func Lax(data []byte, out any, options ...yaml.DecodeOption) error {
	if err := yaml.UnmarshalWithOptions(data, out, options...); err != nil {
		return err
	}

	return nil
}
