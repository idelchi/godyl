package command

import (
	"bytes"
	"context"
	"fmt"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

// Command represents a shell command that can be executed.
type Command string

// String returns the string representation of the Command.
func (c *Command) String() string {
	return string(*c)
}

// From sets the Command value from the provided command string.
func (c *Command) From(command string) {
	*c = Command(command)
}

// Shell executes the Command using mvdan/sh shell interpreter.
// It captures both stdout and stderr output, supports environment variables,
// and returns the stdout output or an error if execution fails.
func (c *Command) Shell(env ...string) (string, error) {
	var stdoutBuf, stderrBuf bytes.Buffer

	// Parse the command string into a shell script
	parser := syntax.NewParser()

	file, err := parser.Parse(bytes.NewReader([]byte(*c)), "")
	if err != nil {
		return "", fmt.Errorf("parsing shell command: %w", err)
	}

	// Set up the interpreter with stdout and stderr buffers to capture output
	runner, err := interp.New(
		interp.StdIO(nil, &stdoutBuf, &stderrBuf),
		interp.Env(expand.ListEnviron(env...)), // Pass the environment variables
	)
	if err != nil {
		return "", fmt.Errorf("setting up shell interpreter: %w", err)
	}

	// Execute the parsed command
	err = runner.Run(context.TODO(), file)
	if err != nil {
		return "", fmt.Errorf(
			"%w: %w: stdout: %s: stderr: %s",
			ErrRun,
			err,
			stdoutBuf.String(),
			stderrBuf.String(),
		)
	}

	// Return the captured stdout output
	return stdoutBuf.String(), nil
}

// ErrRun indicates a failure while executing a shell command.
var ErrRun = fmt.Errorf("running shell command")
