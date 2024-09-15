package debug

import (
	"fmt"
	"os"

	"github.com/kr/pretty"
)

func Debug(format string, args ...any) {
	if os.Getenv("DEBUG") != "" {
		fmt.Printf("DEBUG: "+format+"\n", args...)
	}
}

func Print(a ...any) {
	pretty.Print(a)
}
