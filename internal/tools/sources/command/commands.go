package command

import (
	"fmt"
	"strings"

	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/sources/common"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/file"
)

// Commands represents a slice of shell commands.
//
// TODO(Idelchi): Change to a struct that holds data too, and make a custom unmarshal to still unmarshal into a slice of commands
type Commands []Command

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
// with each command joined by semicolons.
func (c *Commands) Combined() Command {
	stringCommands := make([]string, len(*c))
	for i, cmd := range *c {
		stringCommands[i] = string(cmd)
	}

	return Command(strings.Join(stringCommands, "; "))
}

// Install runs the combined commands for installation using the provided InstallData,
// captures the output, and returns it alongside any errors or found file information.
func (c *Commands) Install(d common.InstallData) (output string, found file.File, err error) {
	cmd := c.Combined()

	// Execute the combined command
	output, err = cmd.Shell(d.Env.ToSlice()...)
	if err != nil {
		return output, "", fmt.Errorf("running combined commands: %w", err)
	}

	return strings.TrimRight(output, "\n"), "", nil
}

func (c *Commands) Exe(env env.Env) (output string, err error) {
	cmd := c.Combined()

	// Execute the combined command
	output, err = cmd.Shell(env.ToSlice()...)
	if err != nil {
		return output, fmt.Errorf("running combined commands: %w", err)
	}

	return strings.TrimRight(output, "\n"), nil
}
