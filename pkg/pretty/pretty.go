// Package pretty contains functions for prettifying and visualizing data.
package pretty

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/showa-93/go-mask"
	"gopkg.in/yaml.v3"
)

func PrintJSON(obj any) {
	fmt.Println(JSON(obj))
}

func PrintYAML(obj any) {
	fmt.Println(YAML(obj))
}

func YAML(obj any) string {
	buf := bytes.Buffer{}
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	err := enc.Encode(&obj)
	if err != nil {
		return err.Error()
	}

	return buf.String()
}

func YAMLMasked(obj any) string {
	return YAML(JSONMasked(obj))
}

// JSON returns a prettified JSON representation of the provided object.
func JSON(obj any) string {
	bytes, err := json.MarshalIndent(obj, "  ", "    ")
	if err != nil {
		return err.Error()
	}

	return string(bytes)
}

// PrintJSONMasked returns a pretty-printed JSON string representation of the provided object with masked sensitive
// fields.
func JSONMasked(obj any) string {
	return JSON(MaskJSON(obj))
}

// JSONMasked returns a pretty-printed JSON representation of the provided object with masked sensitive fields.
func MaskJSON(obj any) any {
	masker := mask.NewMasker()

	masker.SetMaskChar("-")

	masker.RegisterMaskStringFunc(mask.MaskTypeFilled, masker.MaskFilledString)

	t, err := mask.Mask(obj)
	if err != nil {
		return err.Error()
	}

	return t
}
