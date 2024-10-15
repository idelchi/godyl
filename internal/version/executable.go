package version

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/idelchi/godyl/pkg/file"
)

// Executable consists of a full path to a file and its version.
// An attempt to parse the version into a string can be done by using a `Version` type.
type Executable struct {
	File    file.File
	Version string
}

func NewExecutable(paths ...string) Executable {
	return Executable{File: file.New(paths...)}
}

func (e Executable) Command(ctx context.Context, cmdArgs []string) (string, error) {
	var out bytes.Buffer

	cmd := exec.CommandContext(ctx, e.File.Name(), cmdArgs...)
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()

	return strings.TrimSpace(out.String()), err
}

// ParseVersion attempts to parse the version of the executable using the provided Version object.
func (e *Executable) ParseVersion() error {
	timeout := 30 * time.Second

	// Create a context with a timeout to prevent hanging
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	version := NewDefaultVersionParser()

	// Iterate through each command strategy
	for _, cmdArgs := range version.Commands {
		// Get the output of the command
		output, err := e.Command(ctx, cmdArgs)
		if err != nil {
			fmt.Printf("Error parsing version: %v: %q\n", err, output)
			continue
		}

		if version, err := version.ParseString(output); err == nil {
			e.Version = version

			return nil
		} else {
			fmt.Printf("Error parsing version: %v: %q: %q\n", err, version, output)

			continue
		}
	}

	e.Version = ""

	return errors.New("unable to parse version from output")
}
