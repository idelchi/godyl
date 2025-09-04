//go:generate go tool string-enumer -t Match -o match_enumer___generated.go .
//go:generate go tool string-enumer -t Type -o type_enumer___generated.go .
package hints

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/goccy/go-yaml/ast"

	"github.com/idelchi/godyl/pkg/unmarshal"
)

// Type represents different pattern matching types for hints.
type Type string

const (
	// Glob indicates that the hint is a glob pattern.
	Glob Type = "glob"
	// Regex indicates that the hint is a regular expression.
	Regex Type = "regex"
	// GlobStar indicates that the hint is a globstar pattern.
	GlobStar Type = "globstar"
	// StartsWith indicates that the hint is a startswith pattern.
	StartsWith Type = "startswith"
	// EndsWith indicates that the hint is an endswith pattern.
	EndsWith Type = "endswith"
	// Contains indicates that the hint is a contains pattern.
	Contains Type = "contains"
)

// Match represents different matching strategies for hints.
type Match string

const (
	// Weighted indicates that the hint is a weighted match.
	Weighted Match = "weighted"
	// Required indicates that the hint is a required match.
	Required Match = "required"
	// Excluded indicates that the hint is an excluded match.
	Excluded Match = "excluded"
)

// Hint represents a pattern used to match asset names.
// It can be a regular expression or a simple string pattern.
type Hint struct {
	// Pattern to match against the asset's name.
	Pattern string `single:"true"`
	// Weight used to adjust the score for non-mandatory hints.
	Weight unmarshal.Templatable[int]
	// Type determines the engine used to match the pattern.
	Type Type `validate:"oneof=glob regex globstar startswith endswith contains"`
	// Match indicates the type of match.
	Match unmarshal.Templatable[Match]
}

// Matches evaluates whether the hint pattern matches the given string.
func (h Hint) Matches(s string) (match bool, err error) {
	switch h.Type {
	case Glob:
		return path.Match(h.Pattern, s)
	case Regex:
		regex, err := regexp.Compile(h.Pattern)
		if err != nil {
			return false, err
		}

		return regex.MatchString(s), nil
	case GlobStar:
		return doublestar.Match(h.Pattern, s)
	case StartsWith:
		return strings.HasPrefix(s, h.Pattern), nil
	case EndsWith:
		return strings.HasSuffix(s, h.Pattern), nil
	case Contains:
		return strings.Contains(s, h.Pattern), nil
	default:
		return false, fmt.Errorf("unknown hint type: %q", h.Type)
	}
}

// UnmarshalYAML implements the yaml.Unmarshaler interface for Hint.
func (h *Hint) UnmarshalYAML(node ast.Node) error {
	type raw Hint

	if err := unmarshal.SingleStringOrStruct(node, (*raw)(h)); err != nil {
		return fmt.Errorf("unmarshalling hint: %w", err)
	}

	return nil
}

// Parse validates and prepares the hint for matching operations.
func (h *Hint) Parse() (err error) {
	if h.Weight.IsUnset() {
		h.Weight.Set("1")
	}

	if h.Match.IsUnset() {
		h.Match.Set(string(Weighted))
	}

	if h.Type == "" {
		h.Type = Glob
	}

	if err := h.Weight.Parse(); err != nil {
		return fmt.Errorf("parsing weight: %w", err)
	}

	if err := h.Match.Parse(); err != nil {
		return fmt.Errorf("parsing match: %w", err)
	}

	if !h.Match.Value.Valid() {
		return fmt.Errorf("invalid match type: %q", h.Match.Value)
	}

	if !h.Type.Valid() {
		return fmt.Errorf("invalid hint type: %q", h.Type)
	}

	return err
}
