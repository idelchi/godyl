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

// Commands represents a slice of shell commands.
type Commands struct {
	// Commands hold the commands to execute
	Commands unmarshal.SingleOrSliceType[Command]
	// AllowFailure indicates whether the return code of the commands (combined) should be suppressed or not.
	AllowFailure bool `yaml:"allow_failure"`
	// ExitOnError indicates whether to exit on error (injects `set -e -o pipefail` or `set +e` depending on the value).
	ExitOnError bool `yaml:"exit_on_error"`

	// Data holds additional metadata related to the repository.
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

// Get retrieves a specific attribute of the commands.
func (c *Commands) Get(_ string) string {
	return "N/A"
}

// Initialize prepares the Commands based on the given string.
func (c *Commands) Initialize(command string) error {
	return nil
}

// // Exe just satisfies the interface for the Commands struct.
// func (c *Commands) Exe() error {
// 	return nil
// }

// Version just satisfies the interface for the Commands struct.
func (c *Commands) Version(version string) error {
	return nil
}

// Path just satisfies the interface for the Commands struct.
func (c *Commands) Path(path string, patterns []string, _ string, _ match.Requirements) error {
	return nil
}

// Combined returns all commands in the Commands slice as a single Command,
// with commands joined by semicolons. It prepends "set -e" only if AllowFailure is false.
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

func (c *Commands) Exe(env env.Env) (output string, err error) {
	cmd := c.Combined()

	// Execute the combined command
	output, err = cmd.Shell(env.ToSlice()...)
	if err != nil && !(c.AllowFailure && errors.Is(err, ErrRun)) {
		return output, fmt.Errorf("running combined commands: %w", err)
	}

	return strings.TrimRight(output, "\n"), nil
}
