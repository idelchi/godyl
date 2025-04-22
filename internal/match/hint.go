package match

import (
	"fmt"
	"strconv"

	"github.com/fatih/structs"

	"github.com/idelchi/godyl/pkg/unmarshal"

	"gopkg.in/yaml.v3"
)

// Hint represents a pattern used to match asset names.
// It can be a regular expression or a simple string pattern.
type Hint struct {
	Pattern string // Pattern to match against the asset's name.
	Weight  string // Weight used to adjust the score for non-mandatory hints.
	Must    bool   // Indicates if the hint is mandatory for a match.

	weight int `json:"-" mapstructure:"-" yaml:"-"`
}

func (h *Hint) UnmarshalYAML(value *yaml.Node) error {
	h.Weight = "1"

	// Perform custom unmarshaling with field validation, allowing only known fields.
	type raw Hint

	return unmarshal.DecodeWithOptionalKnownFields(value, (*raw)(h), true, structs.New(h).Name())
}

func (h *Hint) Parse() (err error) {
	// Parse the condition string into a boolean value.
	h.weight, err = strconv.Atoi(h.Weight)
	if err != nil {
		return fmt.Errorf("parsing weight %q: %w", h.Weight, err)
	}

	return err
}

// GetWeight returns the weight of the hint.
func (h Hint) GetWeight() int {
	return h.weight
}
