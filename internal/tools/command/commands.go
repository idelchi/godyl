package command

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/structs"

	"github.com/idelchi/godyl/internal/tools/sources/common"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/unmarshal"

	"gopkg.in/yaml.v3"
)

// Commands represents a collection of shell commands that can be executed together.
type Commands struct {
	Data         common.Metadata `yaml:"-"`
	Commands     unmarshal.SingleOrSliceType[Command]
	AllowFailure bool `yaml:"allow_failure"`
	ExitOnError  bool `yaml:"exit_on_error"`
}

func (e *Commands) UnmarshalYAML(value *yaml.Node) error {
	e.ExitOnError = true

	if value.Kind == yaml.ScalarNode {
		e.Commands = []Command{Command(value.Value)}

		return nil
	}

	// Handle sequence node directly (list of commands)
	if value.Kind == yaml.SequenceNode {
		var commands []Command
		if err := value.Decode(&commands); err != nil {
			return err
		}

		e.Commands = commands

		return nil
	}

	// Perform custom unmarshaling with field validation, allowing only known fields.
	type raw Commands

	return unmarshal.DecodeWithOptionalKnownFields(value, (*raw)(e), true, structs.New(e).Name())
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
	output, err = cmd.Shell(env.ToSlice()...)
	if err != nil && !(c.AllowFailure && errors.Is(err, ErrRun)) {
		return output, fmt.Errorf("running combined commands: %w", err)
	}

	return strings.TrimRight(output, "\n"), nil
}
