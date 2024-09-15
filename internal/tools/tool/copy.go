package tool

import (
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/idelchi/godyl/pkg/utils"
)

func (t *Tool) CopyYAML() error {
	bytes, err := yaml.Marshal(*t)
	if err != nil {
		return fmt.Errorf("marshaling %q to yaml: %w", t.Name, err)
	}

	err = json.Unmarshal(bytes, t)
	if err != nil {
		return fmt.Errorf("unmarshaling %q from yaml: %w", t.Name, err)
	}

	return nil
}

func (t *Tool) CopiedYAML() (*Tool, error) {
	bytes, err := yaml.Marshal(*t)
	if err != nil {
		return nil, fmt.Errorf("marshaling %q to yaml: %w", t.Name, err)
	}

	dst := &Tool{}

	err = json.Unmarshal(bytes, &dst)
	if err != nil {
		return dst, fmt.Errorf("unmarshaling %q from yaml: %w", t.Name, err)
	}

	return dst, nil
}

// Copied creates and returns a deep copy of the Tool instance.
func (t *Tool) Copied() (*Tool, error) {
	copied, err := utils.DeepCopy(*t)
	if err != nil {
		return nil, fmt.Errorf("copying tool %q: %w", t.Name, err)
	}

	return &copied, nil
}

// Copied copies the Tool instance and sets it to itself.
func (t *Tool) Copy() error {
	copied, err := utils.DeepCopy(*t)
	if err != nil {
		return fmt.Errorf("copying tool %q: %w", t.Name, err)
	}

	*t = copied

	return nil
}
