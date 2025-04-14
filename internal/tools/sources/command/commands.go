package command

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/structs"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/sources/common"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/unmarshal"
	"gopkg.in/yaml.v3"
)

// Commands represents a collection of shell commands that can be executed together.
type Commands struct {
	// Commands is the list of shell commands to execute.
	Commands unmarshal.SingleOrSliceType[Command]

	// AllowFailure determines if command execution should continue even if a command fails.
	AllowFailure bool `yaml:"allow_failure"`

	// ExitOnError controls shell error handling behavior.
	// When true, injects 'set -e -o pipefail', when false injects 'set +e'.
	ExitOnError bool `yaml:"exit_on_error"`

	// Data contains additional metadata about the command source.
	Data common.Metadata `yaml:"-"`
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

// Get retrieves a specific attribute of the Commands.
// Currently returns a placeholder value as this is an interface requirement.
func (c *Commands) Get(_ string) string {
	return "N/A"
}

// Initialize prepares the Commands based on the given command string.
// Currently a no-op as initialization is handled elsewhere.
func (c *Commands) Initialize(command string) error {
	return nil
}

// // Exe just satisfies the interface for the Commands struct.
// func (c *Commands) Exe() error {
// 	return nil
// }

// Version satisfies the Populater interface requirement.
// Currently a no-op as version handling is not needed for Commands.
func (c *Commands) Version(version string) error {
	return nil
}

// Path satisfies the Populater interface requirement.
// Currently a no-op as path handling is not needed for Commands.
func (c *Commands) Path(path string, patterns []string, _ string, _ match.Requirements) error {
	return nil
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

// Exe executes the combined commands with the provided environment variables.
// Returns the command output and any execution errors, respecting AllowFailure setting.
func (c *Commands) Exe(env env.Env) (output string, err error) {
	cmd := c.Combined()

	// Execute the combined command
	output, err = cmd.Shell(env.ToSlice()...)
	if err != nil && !(c.AllowFailure && errors.Is(err, ErrRun)) {
		return output, fmt.Errorf("running combined commands: %w", err)
	}

	return strings.TrimRight(output, "\n"), nil
}
