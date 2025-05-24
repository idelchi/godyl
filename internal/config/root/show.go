package root

import (
	"github.com/idelchi/godyl/pkg/pretty"
)

// ShowFuncType declares the type for the ShowFunc function.
type ShowFuncType func() func(any)

// Verbosity controls the level of output shown to the user in the ShowFunc function.
type Verbosity int

const (
	// None means no output.
	None Verbosity = iota
	// Masked means masked output.
	Masked
	// Default means unmasked output.
	Default
)

// ShowFunc returns the function to be used for showing the output based on the Verbosity level.
func (r *Root) ShowFunc() func(any) {
	switch r.Show {
	case None:
		return nil
	case Masked:
		return pretty.PrintYAMLMasked
	case Default:
		return pretty.PrintYAML
	default:
		return pretty.PrintYAML
	}
}
