// Package mode provides functionality for specifying tool operational modes.
package mode

// Mode represents the operational mode for a tool, determining how the tool is handled.
// It could specify different behaviors such as extracting the tool from an archive or searching for the tool.
type Mode string

const (
	// Extract mode indicates that the tool should be extracted from an archive or compressed file.
	Extract Mode = "extract"
	// Find mode indicates that the tool should be located within a specified directory or environment.
	Find Mode = "find"
)

// String returns the string representation of the Mode.
func (m Mode) String() string {
	return string(m)
}
