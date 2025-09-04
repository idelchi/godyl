package command

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml/ast"

	"github.com/idelchi/godyl/internal/tools/sources/install"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/unmarshal"
)

// Commands represents a collection of shell commands that can be executed together.
type Commands struct {
	Data         install.Metadata                     `yaml:"-"`
	Commands     unmarshal.SingleOrSliceType[Command] `yaml:"commands"`
	AllowFailure bool                                 `yaml:"allow-failure"`
	ExitOnError  bool                                 `yaml:"exit-on-error"`
}

// UnmarshalYAML implements custom YAML unmarshaling for Commands.
// It supports both single command strings and command arrays.
func (e *Commands) UnmarshalYAML(node ast.Node) error {
	e.ExitOnError = true

	switch node.(type) {
	// When given on the style
	// commands: "echo hello world"
	// or
	// commands:
	//   - echo "hello"
	//   - echo "world"
	case *ast.StringNode, *ast.SequenceNode:
		result, err := unmarshal.SingleOrSlice[Command](node)
		if err != nil {
			return fmt.Errorf("unmarshaling commands: %w", err)
		}

		e.Commands = result

		return nil
	}

	type raw Commands

	return unmarshal.Decode(node, (*raw)(e))
}

// Combined joins all commands into a single Command with proper shell options.
// Prepends shell error handling options based on ExitOnError setting and
// joins commands with semicolons.
func (e *Commands) Combined() Command {
	stringCommands := make([]string, 0, len(e.Commands)+1)

	if e.ExitOnError {
		stringCommands = append(stringCommands, "set -e -o pipefail")
	} else {
		stringCommands = append(stringCommands, "set +e")
	}

	for _, cmd := range e.Commands {
		stringCommands = append(stringCommands, string(cmd))
	}

	return Command(strings.Join(stringCommands, "\n"))
}

// Run executes the combined commands with the provided environment variables.
// Returns the command output and any execution errors, respecting AllowFailure setting.
func (e *Commands) Run(ctx context.Context, env env.Env) (output string, err error) {
	cmd := e.Combined()

	// Execute the combined command
	output, err = cmd.Shell(ctx, env.AsSlice()...)
	if err != nil && (!e.AllowFailure || !errors.Is(err, ErrRun)) {
		return output, fmt.Errorf("running combined commands: %w", err)
	}

	return strings.TrimRight(output, "\n"), nil
}
