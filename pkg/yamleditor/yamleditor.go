package yamleditor

import (
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"

	"dario.cat/mergo"
)

// Editor provides in-place, comment-preserving YAML editing.
type Editor struct {
	data map[string]any
}

// New builds an Editor from raw YAML bytes (comments are discarded).
func New(src []byte) (*Editor, error) {
	var m map[string]any
	if err := yaml.Unmarshal(src, &m); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return &Editor{data: m}, nil
}

// Set assigns v at dot-separated path, creating parents as needed.
func (e *Editor) Set(path string, v any) error {
	overlay := buildNestedMap(strings.Split(path, "."), v)

	return mergo.Merge(&e.data, overlay, mergo.WithOverride)
}

// SetRaw parses yamlFragment and delegates to Set.
func (e *Editor) SetRaw(path, yamlFragment string) error {
	var val any
	if err := yaml.Unmarshal([]byte(yamlFragment), &val); err != nil {
		return fmt.Errorf("SetRaw unmarshal: %w", err)
	}

	return e.Set(path, val)
}

// Get retrieves the value at path.
func (e *Editor) Get(path string) (any, error) {
	cur := any(e.data)
	for k := range strings.SplitSeq(path, ".") {
		mp, ok := cur.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("%q: non-map segment", k)
		}

		var exists bool

		cur, exists = mp[k]
		if !exists {
			return nil, fmt.Errorf("%q: key not found", k)
		}
	}

	return cur, nil
}

// Render returns the YAML encoding of the current document.
func (e *Editor) Render() ([]byte, error) {
	return yaml.Marshal(e.data)
}

func buildNestedMap(keys []string, v any) map[string]any {
	if len(keys) == 0 {
		return nil
	}

	if len(keys) == 1 {
		return map[string]any{keys[0]: v}
	}

	return map[string]any{
		keys[0]: buildNestedMap(keys[1:], v),
	}
}
