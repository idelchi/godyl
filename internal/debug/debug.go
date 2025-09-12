// Package debug provides debugging utilities for development and troubleshooting.
package debug

import (
	"fmt"
	"os"

	"github.com/kr/pretty"
)

// Debug prints formatted debug messages when GODYL_DEBUG environment variable is set.
func Debug(format string, args ...any) {
	if os.Getenv("GODYL_DEBUG") != "" {
		fmt.Printf( //nolint:forbidigo // Debug package is meant for printing
			"DEBUG: "+format+"\n",
			args...)
	}
}

// Print outputs debug information using pretty printing when DEBUG environment variable is set.
func Print(a ...any) {
	pretty.Print(a) //nolint:gosec,errcheck // Debug print functions do not need error handling
}
