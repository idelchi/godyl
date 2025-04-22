package executable

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/idelchi/godyl/pkg/path/file"
)

// Executable represents a command-line executable.
type Executable struct {
	file.File
}

// NewExecutable creates a new Executable instance from the provided paths.
func New(paths ...string) Executable {
	return Executable{File: file.New(paths...)}
}

// Command runs the specified command arguments by passing them to the executable.
// It returns the output of the command as a trimmed string and any error encountered during execution.
func (e Executable) Command(ctx context.Context, cmdArgs []string) (string, error) {
	var out bytes.Buffer

	cmd := exec.CommandContext(ctx, e.Path(), cmdArgs...)
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Stdin = os.Stdin

	err := cmd.Run()

	return strings.TrimSpace(out.String()), err
}

// Parse attempts to parse the output of the executable using the provided Parser object.
// It iterates over each command defined in the Parser and returns the first successful match.
func (e Executable) Parse(parser *Parser) (string, error) {
	const timeout = 60 * time.Second

	// Create a context with a timeout to prevent hanging
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	errs := make([]error, 0, len(parser.Commands))

	// Iterate through each command strategy
	for _, cmdArgs := range parser.Commands {
		// Get the output of the command
		output, err := e.Command(ctx, []string{cmdArgs})
		if err != nil {
			// // Collect errors and continue
			errs = append(errs, err)
		}

		if output, err := parser.Parse(output); err == nil {
			return output, nil
		}

		errs = append(errs, err)

		continue
	}

	// Join all errors into a single error message
	return "", fmt.Errorf("parsing output: %w", errors.Join(errs...))
}
