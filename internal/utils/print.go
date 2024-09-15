package utils

import (
	"fmt"

	"github.com/idelchi/godyl/pkg/pretty"
)

// Print displays the configuration in the specified format.
func Print(format string, cfg ...any) {
	printFunc := func(any) {}

	switch format {
	case "json":
		printFunc = pretty.PrintJSONMasked
	case "yaml":
		printFunc = pretty.PrintYAMLMasked
	case "env":
		printFunc = pretty.PrintEnv
	default:
		fmt.Printf("unsupported output format: %s\n", format)
	}

	for _, c := range cfg {
		printFunc(c)
	}
}
