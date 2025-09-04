package goi

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
)

// Installer handles the installation of Go binaries using the provided Binary.
type Installer struct {
	Binary Binary // Binary represents the Go binary used for the installation process.
	GOOS   string // Target operating system for cross-compilation.
	GOARCH string // Target architecture for cross-compilation.
}

// Install executes the `go install` command for the provided package path.
// It captures both stdout and stderr, returning them as output, and reports errors if the installation fails.
func (i *Installer) Install(path string) (output string, err error) {
	var stdoutBuf, stderrBuf bytes.Buffer

	// Prepare the command
	cmd := exec.CommandContext(context.Background(), //nolint:gosec // Path is validated by upstream tool configuration
		i.Binary.File.Path(),
		"install",
		path,
	)

	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, i.Binary.Env.ToSlice()...)

	if i.GOOS != "" {
		cmd.Env = append(cmd.Env, "GOOS="+i.GOOS)
	}

	if i.GOARCH != "" {
		cmd.Env = append(cmd.Env, "GOARCH="+i.GOARCH)
	}

	// Capture stdout and stderr
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	// Run the command
	if err := cmd.Run(); err != nil {
		return stdoutBuf.String() + "\n" + stderrBuf.String(), fmt.Errorf(
			"go install: %w: %s",
			err,
			stdoutBuf.String()+"\n"+stderrBuf.String(),
		)
	}

	// Return both stdout and stderr as the output
	return stdoutBuf.String() + "\n" + stderrBuf.String(), nil
}
