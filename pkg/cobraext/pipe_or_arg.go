package cobraext

import (
	"fmt"

	"github.com/idelchi/gogen/pkg/stdin"
)

// 3. Empty string if neither source available.
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
