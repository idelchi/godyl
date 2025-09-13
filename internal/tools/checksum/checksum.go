// Package checksum provides a structure for defining and working with checksums.
package checksum

import (
	"strings"

	"github.com/goccy/go-yaml/ast"

	"github.com/idelchi/godyl/internal/debug"
	"github.com/idelchi/godyl/pkg/unmarshal"
)

// Checksum represents a checksum configuration with type, value, and optionality.
type Checksum struct {
	// Type is the type of checksum (e.g., sha256, sha512, sha1, md5, file, none).
	Type string `validate:"oneof=sha256 sha512 sha1 md5 file none" single:"true"`
	// Value is the checksum value or URL.
	Value string
	// Optional indicates if the checksum is optional.
	Optional bool
}

// UnmarshalYAML implements custom YAML unmarshaling for Checksum configuration.
// Supports both scalar values (treated as Type) and map values.
func (c *Checksum) UnmarshalYAML(node ast.Node) error {
	type raw Checksum

	return unmarshal.SingleStringOrStruct(node, (*raw)(c))
}

// IsSet returns true if the checksum has a value defined.
func (c *Checksum) IsSet() bool {
	return c.Value != ""
}

// Indicators returns a list of common substrings found in checksum file names.
func Indicators() []string {
	return []string{
		"checksum",
		"checksums",
		"sha256",
		"sha-256",
		"sha512",
		"sha-512",
		"md5",
		"md-5",
		"SHASUMS",
	}
}

// IsChecksumLike determines if a name is likely a checksum file.
func IsChecksumLike(name string) bool {
	for _, indicator := range Indicators() {
		debug.Debug("checking if %q contains %q", name, indicator)
		if strings.Contains(name, indicator) {
			return true
		}
	}

	return false
}
