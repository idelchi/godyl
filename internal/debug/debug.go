// Package debug provides debugging utilities for development and troubleshooting.
package debug

import (
	"fmt"
	"os"
)

// Debug prints formatted debug messages when GODYL_DEBUG environment variable is set to "true".
func Debug(format string, args ...any) {
	if os.Getenv("GODYL_DEBUG") == "true" {
		fmt.Printf( //nolint:forbidigo // Debug package is meant for printing
			"DEBUG: "+format+"\n",
			args...)
	}
}
