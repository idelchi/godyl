// Package checksum provides a structure for defining and working with checksums.
package checksum

import (
	"errors"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml/ast"

	"github.com/idelchi/godyl/internal/data"
	"github.com/idelchi/godyl/pkg/download"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/unmarshal"
)

// Checksum represents a checksum configuration with type, value, and optionality.
type Checksum struct {
	// Type is the type of checksum (e.g., sha256, sha512, sha1, md5, file, none).
	Type string `validate:"oneof=sha256 sha512 sha1 md5 file none" single:"true"`
	// Value as a checksum string or a URL/path to a file containing the checksum.
	Value string
	// Pattern is an optional glob pattern to consider when selecting the checksum file with the combination
	// `Type: file` and `Value: ""`. It is ignored for other types.
	Pattern string
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

// IsMandatory returns true if the checksum is mandatory.
func (c *Checksum) IsMandatory() bool {
	return c.Type != "none"
}

// Resolve determines the value type if type is not none or file, and the value
// looks url/path-like.
func (c *Checksum) Resolve(skipVerifySSL bool) error {
	// For none and file type, do nothing
	if c.Type == "none" {
		return nil
	}

	// If the value starts with sha256 sha512 sha1 md5, strip it and set the type and value
	for _, algo := range []string{"sha256", "sha512", "sha1", "md5"} {
		if value, ok := strings.CutPrefix(c.Value, algo+":"); ok {
			c.Type = algo
			c.Value = strings.TrimSpace(value)

			return nil
		}
	}

	// For file type, do nothing
	if c.Type == "file" {
		return nil
	}

	if url, ok := strings.CutPrefix(c.Value, "url:"); ok {
		options := []download.Option{}

		if skipVerifySSL {
			options = append(options, download.WithInsecureSkipVerify())
		}

		dir, err := data.CreateUniqueDirIn()
		if err != nil {
			return fmt.Errorf("creating random dir: %w", err)
		}

		defer func() {
			err = errors.Join(err, dir.Remove())
		}()

		checksum, err := download.New(options...).Download(url, dir.Path())
		if err != nil {
			return fmt.Errorf("downloading checksum from %q: %w", url, err)
		}

		bytes, err := checksum.Read()
		if err != nil {
			return fmt.Errorf("reading checksum file from path %q: %w", dir.Path(), err)
		}

		c.Value = strings.TrimSpace(string(bytes))

		// If c.Value contains spaces, the first part is the checksum type
		c.Value, _, _ = strings.Cut(c.Value, " ")

		return nil
	}

	if path, ok := strings.CutPrefix(c.Value, "path:"); ok {
		bytes, err := file.New(path).Read()
		if err != nil {
			return fmt.Errorf("reading checksum file from path %q: %w", path, err)
		}

		c.Value = strings.TrimSpace(string(bytes))

		// If c.Value contains spaces, the first part is the checksum type
		c.Value, _, _ = strings.Cut(c.Value, " ")

		return nil
	}

	return nil
}

// ToQuery converts the checksum to a query string format.
func (c *Checksum) ToQuery() string {
	return "checksum=" + c.Type + ":" + c.Value
}
