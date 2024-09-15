package executable

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Parser provides functionality for parsing strings from an output using regex patterns.
type Parser struct {
	// Patterns are regex patterns used to match strings in command output.
	Patterns []string
	// Commands are command strategies to be executed for extraction.
	Commands []string
}

// Parse attempts to extract the version string from the provided output string using the defined regex patterns.
// It normalizes multi-line output into a single line and tries to match the patterns.
// Returns the first matched version string or an error if no match is found.
func (p *Parser) Parse(output string) (string, error) {
	// Normalize the output into a single line by replacing newlines with spaces.
	normalizedOutput := strings.Join(strings.Split(output, "\n"), " ")

	// Try to match each regex pattern on the normalized output.
	for _, pattern := range p.Patterns {
		p := regexp.MustCompile(pattern)

		if matches := p.FindStringSubmatch(normalizedOutput); len(matches) > 1 {
			// Return the first matched version group.
			return matches[1], nil
		}
	}

	return "", fmt.Errorf("%w: %v", ErrNoMatch, p.Patterns)
}

// ErrNoMatch is returned when no match is found in the output.
var ErrNoMatch = errors.New("no match found")
