package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/sources/common"
	"github.com/idelchi/godyl/pkg/file"
)

// Commands represents a slice of shell commands.
type Commands []Command

// Get retrieves a specific attribute of the commands.
func (c *Commands) Get(_ string) string {
	if len(*c) == 0 {
		return ""
	}
	// Return the first command as a string representation
	return string((*c)[0])
}

// Initialize prepares the Commands based on the given string.
func (c *Commands) Initialize(command string) error {
	if command == "" {
		return nil
	}

	*c = append(*c, Command(command))

	return nil
}

// Exe executes the commands.
func (c *Commands) Exe() error {
	if len(*c) == 0 {
		return nil
	}

	// Execute all commands in sequence
	for _, cmd := range *c {
		// Get current environment variables
		env := os.Environ()

		// Execute the command with the environment
		if _, err := cmd.Shell(env...); err != nil {
			return fmt.Errorf("executing command: %w", err)
		}
	}

	return nil
}

// Version sets the version for the commands.
func (c *Commands) Version(version string) error {
	if version == "" {
		return nil
	}
	// Store version as a command
	*c = append(*c, Command(version))

	return nil
}

// Path sets up the path for the commands, using the provided parameters.
func (c *Commands) Path(path string, patterns []string, _ string, _ match.Requirements) error {
	if path == "" {
		return nil
	}
	// Create a command that includes path information
	cmd := Command(fmt.Sprintf("cd %s && %s", path, strings.Join(patterns, " ")))
	*c = append(*c, cmd)

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
