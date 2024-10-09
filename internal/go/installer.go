package ginstaller

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type GInstaller struct {
	Binary Binary
}

func (i *GInstaller) Install(path string) (output string, err error) {
	var stdoutBuf, stderrBuf bytes.Buffer

	// Prepare the command
	cmd := exec.Command(i.Binary.Path.Path, "install", path)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, i.Binary.Env.ToSlice()...)

	// Capture stdout and stderr
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	// Run the command
	if err := cmd.Run(); err != nil {
		return stdoutBuf.String() + "\n" + stderrBuf.String(), fmt.Errorf("go install: %w", err)
	}

	// Return both stdout and stderr as the output
	return stdoutBuf.String() + "\n" + stderrBuf.String(), nil
}
