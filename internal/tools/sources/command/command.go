package command

import (
	"bytes"
	"context"
	"fmt"

	"mvdan.cc/sh/v3/expand"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
)

type Command string

func (c Command) String() string {
	return string(c)
}

func (c *Command) From(command string) {
	*c = Command(command)
}

// Run executes the command using mvdan/sh, capturing output and returning it.
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

	// Return the captured output
	return stdoutBuf.String(), nil
}
