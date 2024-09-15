package unmarshal

import (
	"github.com/goccy/go-yaml"
)

func Strict(data []byte, out any) error {
	if err := yaml.UnmarshalWithOptions(data, out, yaml.Strict()); err != nil {
		return err
	}

	return nil
}
