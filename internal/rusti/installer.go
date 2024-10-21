package rusti

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type Installer struct {
	Binary Binary
}

func (i *Installer) Install(path string) (output string, err error) {
	var stdoutBuf, stderrBuf bytes.Buffer

	cargoPath := i.Binary.File.Dir().Join("cargo").String()
	cmd := exec.Command(cargoPath, "install", path)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, i.Binary.Env.ToSlice()...)

	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	if err := cmd.Run(); err != nil {
		return stdoutBuf.String() + "\n" + stderrBuf.String(), fmt.Errorf("cargo install: %w", err)
	}

	return stdoutBuf.String() + "\n" + stderrBuf.String(), nil
}
