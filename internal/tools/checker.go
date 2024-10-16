package tools

import (
	"github.com/idelchi/godyl/internal/tools/sources/command"
)

// Checker represents a tool checker configuration.
type Checker struct {
	// Test defines the commands that should be run to test or validate the tool's functionality.
	Test command.Commands
	// Checksum is used for verifying the integrity of the tool (e.g., via a hash).
	Checksum Checksum
}

// Checksum represents a checksum configuration.
type Checksum struct {
	// Path to the checksum file.
	Path string
	// Enabled specifies whether checksum verification is enabled.
	Enabled bool
}
