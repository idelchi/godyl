// Package pretty contains functions for prettifying and visualizing data.
package pretty

import (
	"encoding/json"
	"fmt"

	"github.com/showa-93/go-mask"
)

func PrintJSON(obj any) {
	fmt.Println(JSON(obj))
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
func PrintJSONMasked(obj any) string {
	return JSON(JSONMasked(obj))
}

// JSONMasked returns a pretty-printed JSON representation of the provided object with masked sensitive fields.
func JSONMasked(obj any) any {
	masker := mask.NewMasker()

	masker.SetMaskChar("-")

	masker.RegisterMaskStringFunc(mask.MaskTypeFilled, masker.MaskFilledString)

	t, err := mask.Mask(obj)
	if err != nil {
		return err.Error()
	}

	return t
}
