package command

import (
	"errors"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml/ast"

	"github.com/idelchi/godyl/internal/tools/sources/common"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/unmarshal"
)

// Commands represents a collection of shell commands that can be executed together.
type Commands struct {
	Data         common.Metadata                      `json:"-"`
	Commands     unmarshal.SingleOrSliceType[Command] `json:"commands"`
	AllowFailure bool                                 `json:"allow-failure"`
	ExitOnError  bool                                 `json:"exit-on-error"`
}

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
func (c *Commands) Combined() Command {
	stringCommands := make([]string, 0, len(c.Commands)+1)

	if c.ExitOnError {
		stringCommands = append(stringCommands, "set -e -o pipefail")
	} else {
		stringCommands = append(stringCommands, "set +e")
	}

	for _, cmd := range c.Commands {
		stringCommands = append(stringCommands, string(cmd))
	}

	return Command(strings.Join(stringCommands, "; "))
}

// Run executes the combined commands with the provided environment variables.
// Returns the command output and any execution errors, respecting AllowFailure setting.
func (c *Commands) Run(env env.Env) (output string, err error) {
	cmd := c.Combined()

	// Execute the combined command
	output, err = cmd.Shell(env.AsSlice()...)
	if err != nil && (!c.AllowFailure || !errors.Is(err, ErrRun)) {
		return output, fmt.Errorf("running combined commands: %w", err)
	}

	return strings.TrimRight(output, "\n"), nil
}
