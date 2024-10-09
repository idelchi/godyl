package executable

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// TODO(Idelchi): Allow for custom regex patterns and command strategies, passed in the YAML.

// Version type holds version parsing information.
type Version struct {
	Patterns []*regexp.Regexp // List of regex patterns for parsing
	Commands [][]string       // List of version command strategies
	String   string           // String representation of the version
}

// NewDefaultVersion creates a Version object with the simplified regex patterns.
func NewDefaultVersionParser() *Version {
	return &Version{
		Patterns: []*regexp.Regexp{
			// Pattern for X.X.X, surrounded by any characters
			regexp.MustCompile(`.*?(\d+\.\d+\.\d+).*`),
			// Pattern for X.X, surrounded by any characters
			regexp.MustCompile(`.*?(\d+\.\d+).*`),
		},
		Commands: [][]string{
			{"--version"}, // Default attempt with --version
			{"version"},   // Default attempt with version
		},
	}
}

// commandWithContext runs the executable with the provided arguments using a timeout.
func (f Executable) Command(ctx context.Context, cmdArgs []string) (string, error) {
	var out bytes.Buffer
	cmd := exec.CommandContext(ctx, f.Path, cmdArgs...)
	cmd.Stdout = &out
	cmd.Stderr = &out

	return strings.TrimSpace(out.String()), cmd.Run()
}

// ParseString attempts to parse the provided string using the Version patterns.
func (v *Version) ParseString(output string) (string, error) {
	// Normalize the output into a single string (if multi-line)
	normalizedOutput := strings.Join(strings.Split(output, "\n"), " ")

	// Try to match each regex pattern on the whole output
	for _, pattern := range v.Patterns {
		if matches := pattern.FindStringSubmatch(normalizedOutput); len(matches) > 1 {
			// Return the first matched version group from the whole output
			return matches[1], nil
		}
	}

	return "", errors.New("unable to parse version from output")
}

// ParseVersion attempts to parse the version of the executable using the provided Version object.
func (f *Executable) ParseVersion() error {
	timeout := 5 * time.Second

	// Create a context with a timeout to prevent hanging
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	version := NewDefaultVersionParser()

	// Iterate through each command strategy
	for _, cmdArgs := range version.Commands {
		// Get the output of the command
		output, err := f.Command(ctx, cmdArgs)
		if err != nil {
			continue
		}

		// Use the new ParseString method to parse the output
		if version, err := version.ParseString(output); err == nil {
			f.Version = version

			return nil
		} else {
			f.Version = ""

			return nil
		}
	}

	return errors.New("unable to parse version from output")
}
