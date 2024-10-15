package version

import (
	"errors"
	"regexp"
	"strings"
)

// TODO(Idelchi): Allow for custom regex patterns and command strategies, passed in the YAML.

type Version struct {
	Patterns []*regexp.Regexp // List of regex patterns for parsing
	Commands [][]string       // List of version command strategies
	String   string           // String representation of the version
}

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
			{"-version"},  // Default attempt with -version
			{"-v"},        // Default attempt with -v
		},
	}
}

// commandWithContext runs the executable with the provided arguments using a timeout.

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
