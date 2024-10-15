package sources

import (
	"fmt"
	"strings"

	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/pkg/file"
)

type Commands []Command

func (c *Commands) Get(attribute string) string {
	return ""
}

func (*Commands) Initialize(_ string) error {
	return nil
}

func (*Commands) Exe() error {
	return nil
}

func (*Commands) Version(_ string) error {
	return nil
}

func (*Commands) Path(_ string, _ []string, _ string, _ match.Requirements) error {
	return nil
}

func (c *Commands) Install2(_ InstallData) (output, found string, err error) {
	for _, command := range *c {
		output, err := command.Shell()
		output += output + "\n"
		if err != nil {
			return strings.TrimRight(output, "\n"), "", fmt.Errorf("running commands: %w", err)
		}
	}

	return strings.TrimRight(output, "\n"), "", nil
}

// Combined returns all commands as a single Command, joined by semicolons
func (c Commands) Combined() Command {
	stringCommands := make([]string, len(c))
	for i, cmd := range c {
		stringCommands[i] = string(cmd)
	}
	return Command(strings.Join(stringCommands, "; "))
}

func (c Commands) Install(d InstallData) (output string, found file.File, err error) {
	cmd := c.Combined()

	// Execute the combined command
	output, err = cmd.Shell(d.Env.ToSlice()...)
	if err != nil {
		return output, "", fmt.Errorf("running combined commands: %w", err)
	}

	return strings.TrimRight(output, "\n"), "", nil
}
