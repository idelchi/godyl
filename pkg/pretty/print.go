// Package pretty contains functions for prettifying and visualizing data in JSON and YAML formats.
// It includes support for masking sensitive fields when outputting data.

package pretty

import (
	"fmt"
)

// PrintJSON prints a prettified JSON representation of the provided object.
func PrintJSON(obj any) {
	fmt.Println(JSON(obj))
}

// PrintYAML prints a prettified YAML representation of the provided object.
func PrintYAML(obj any) {
	fmt.Println(YAML(obj))
}

// PrintJSONMasked prints a prettified JSON representation of the provided object
// with masked sensitive fields. It uses MaskJSON internally to mask the fields.
func PrintJSONMasked(obj any) {
	fmt.Println(JSONMasked(obj))
}

// PrintYAMLMasked prints a prettified YAML representation of the provided object
// with masked sensitive fields. It uses MaskYAML internally to mask the fields.
func PrintYAMLMasked(obj any) {
	fmt.Println(YAMLMasked(obj))
}
