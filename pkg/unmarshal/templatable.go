package unmarshal

import (
	"errors"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
)

// Templatable wraps a string-based template and its parsed value of type T.
// Used to delay parsing of the value until explicitly needed.
type Templatable[T any] struct {
	// The parsed value of type T.
	Value T
	// The template string to be parsed.
	Template string

	parsed bool
}

// UnmarshalYAML decodes the raw YAML node into the Template string.
// Does not parse the value; use Parse separately for that.
func (t *Templatable[T]) UnmarshalYAML(node ast.Node) error {
	if err := Decode(node, &t.Template); err != nil {
		return fmt.Errorf("unmarshaling template: %w", err)
	}

	return nil
}

// MarshalYAML returns the Template string for YAML serialization.
//
// TODO(Idelchi): Return t.Template or t.Value?
func (t Templatable[T]) MarshalYAML() (any, error) {
	return t.Template, nil
}

// IsUnset returns true if the Template is empty.
func (t *Templatable[T]) IsUnset() bool {
	return t.Template == ""
}

// Set assigns a new value to Template.
func (t *Templatable[T]) Set(s string) {
	t.Template = s
}

// Get returns the parsed value. Does not trigger parsing.
func (t *Templatable[T]) Get() (T, error) {
	if !t.parsed {
		return t.Value, errors.New("value not parsed yet")
	}

	return t.Value, nil
}

// Parse attempts to unmarshal the Template string into Value.
func (t *Templatable[T]) Parse() error {
	if err := yaml.Unmarshal([]byte(t.Template), &t.Value); err != nil {
		return fmt.Errorf("parsing template %q: %w", t.Template, err)
	}

	t.parsed = true

	return nil
}
