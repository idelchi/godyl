package sources

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/idelchi/godyl/internal/match"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

type Command string

// Run executes the command using mvdan/sh, capturing output and returning it
func (c Command) Shell() (string, error) {
	var stdoutBuf, stderrBuf bytes.Buffer

	// Parse the command string into a shell script
	parser := syntax.NewParser()
	file, err := parser.Parse(bytes.NewReader([]byte(c)), "")
	if err != nil {
		return "", fmt.Errorf("parsing shell command: %w", err)
	}

	// Set up the interpreter with stdout and stderr buffers to capture output
	runner, err := interp.New(interp.StdIO(nil, &stdoutBuf, &stderrBuf))
	if err != nil {
		return "", fmt.Errorf("setting up shell interpreter: %w", err)
	}

	// Execute the parsed command
	err = runner.Run(context.TODO(), file)
	if err != nil {
		return "", fmt.Errorf("shell: %w: stderr: %s", err, stderrBuf.String())
	}

	// Return the captured output
	return stdoutBuf.String(), nil
}

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

func (c *Commands) Install(d InstallData) (output string, err error) {
	for _, command := range *c {
		output, err := command.Shell()
		output += output + "\n"
		if err != nil {
			return strings.TrimRight(output, "\n"), fmt.Errorf("running commands: %w", err)
		}
	}

	return strings.TrimRight(output, "\n"), nil
}
