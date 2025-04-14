package cobraext

import (
	"fmt"

	"github.com/idelchi/gogen/pkg/stdin"
)

// PipeOrArg flexibly reads command input from args or pipe.
// Implements a common CLI pattern where input can come from either
// command-line arguments or piped stdin. Priority order:
// 1. First command-line argument if present
// 2. Stdin content if piped
// 3. Empty string if neither source available
func PipeOrArg(args []string) (string, error) {
	if len(args) > 0 {
		// Prioritize argument if it exists, regardless of stdin
		return args[0], nil
	}

	if stdin.IsPiped() {
		// No arg but stdin is piped
		arg, err := stdin.Read()
		if err != nil {
			return "", fmt.Errorf("reading from stdin: %w", err)
		}

		return arg, nil
	}

	return "", nil
}
