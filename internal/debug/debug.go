package debug

import (
	"fmt"
	"os"

	"github.com/kr/pretty"
)

func Debug(format string, args ...any) {
	if os.Getenv("DEBUG") != "" {
		fmt.Printf( //nolint:forbidigo // Debug package is meant for printing
			"DEBUG: "+format+"\n",
			args...)
	}
}

func Print(a ...any) {
	pretty.Print(a) //nolint:errcheck // Debug print functions do not need error handling
}
