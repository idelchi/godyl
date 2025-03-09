package cobraext

import (
	"fmt"

	"github.com/idelchi/gogen/pkg/stdin"
)

// PipeOrArg reads from either the first argument or stdin.
// If an argument is provided, it is returned.
// If no argument is provided but stdin is piped, the stdin content is returned.
// If neither an argument nor stdin is provided, an empty string is returned.
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
