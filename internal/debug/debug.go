package debug

import (
	"fmt"
	"os"
)

func Debug(format string, args ...any) {
	if os.Getenv("DEBUG") != "" {
		fmt.Printf("DEBUG: "+format+"\n", args...)
	}
}
