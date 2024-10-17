package command

import (
	"bytes"
	"context"
	"fmt"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

// Command represents a shell command as a string.
type Command string

// String returns the Command as a string.
func (c Command) String() string {
	return string(c)
}

// From assigns a new shell command string to the Command.
func (c *Command) From(command string) {
	*c = Command(command)
}

// Shell runs the Command using mvdan/sh, capturing both stdout and stderr output.
// It accepts optional environment variables and returns the stdout output and any errors encountered.
func (c Command) Shell(env ...string) (string, error) {
	var stdoutBuf, stderrBuf bytes.Buffer

	// Parse the command string into a shell script
	parser := syntax.NewParser()
	file, err := parser.Parse(bytes.NewReader([]byte(c)), "")
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
			"running shell command: %w: stdout: %s: stderr: %s",
			err,
			stdoutBuf.String(),
			stderrBuf.String(),
		)
	}

	// Return the captured stdout output
	return stdoutBuf.String(), nil
}
