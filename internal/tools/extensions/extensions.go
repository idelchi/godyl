// Package extensions provides functionality for managing tool file extensions.
package extensions

import (
	"fmt"
	"slices"
	"strings"

	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/pkg/unmarshal"

	"gopkg.in/yaml.v3"
)

// Extensions represents a collection of file extensions.
type Extensions unmarshal.SingleOrSliceType[string]

// UnmarshalYAML allows unmarshaling the YAML node as either a single string or a slice of strings.
func (e *Extensions) UnmarshalYAML(value *yaml.Node) error {
	result, err := unmarshal.SingleOrSlice[string](value, false)
	if err != nil {
		return err
	}

	*e = result

	return nil
}

func (e Extensions) Compacted() Extensions {
	return slices.Compact(e)
}

// ToHint converts a list of extensions into a match.Hint.
func (e Extensions) ToHint() match.Hint {
	var noExtensionPart string

	var extensionParts []string

	for _, ext := range e.Compacted() {
		if ext == "" {
			noExtensionPart = "^[^.]+$"
		} else {
			escapedExt := strings.ReplaceAll(ext, ".", `\.`) // Escape dots in extensions
			extensionParts = append(extensionParts, fmt.Sprintf(".*%s$", escapedExt))
		}
	}

	var pattern string
	// Combine both parts
	if noExtensionPart != "" && len(extensionParts) > 0 {
		pattern = fmt.Sprintf("(%s|%s)", noExtensionPart, strings.Join(extensionParts, "|"))
	} else if noExtensionPart != "" {
		pattern = noExtensionPart
	} else {
		pattern = strings.Join(extensionParts, "|")
	}

	return match.Hint{
		Pattern: pattern,
		Must:    true,
	}
}
