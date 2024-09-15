//go:generate go tool string-enumer -t Format -o format_enumer___generated.go .
package iutils

import (
	"github.com/idelchi/godyl/pkg/pretty"
)

// Format represents the output format for displaying configuration.
type Format string

const (
	// JSON represents JSON format.
	JSON Format = "json"
	// YAML represents YAML format.
	YAML Format = "yaml"
	// ENV represents environment variable format.
	ENV Format = "env"
)

// printFunc is a function type for printing configuration.
type printFunc func(any)

// Print displays the configuration in the specified format.
func Print(format Format, cfg ...any) {
	var printF printFunc

	switch format {
	case JSON:
		printF = pretty.PrintJSONMasked
	case YAML:
		printF = pretty.PrintYAMLMasked
	case ENV:
		printF = pretty.PrintEnv
	default:
		printF = pretty.PrintDefault
	}

	for _, c := range cfg {
		printF(c)
	}
}
