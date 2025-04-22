package utils

import (
	"errors"
	"fmt"

	"github.com/fatih/structs"
	"github.com/idelchi/go-next-tag/pkg/stdin"
	"github.com/idelchi/godyl/internal/defaults"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/tools/inherit"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/unmarshal"

	"gopkg.in/yaml.v3"
)

// ToolsLoader is a struct that handles loading tools configuration.
type ToolsLoader struct {
	// File represents the file path to the tools configuration.
	File file.File
	// Defaults represents the default tool configuration.
	Defaults *defaults.Defaults
	// tools is a collection of tools loaded from the configuration.
	tools tools.Tools
}

// LoadTools loads the tools configuration from the specified file path,
// or from stdin if the path is "-". It returns a Tools collection,
// initialized with the provided defaults.
func LoadTools(path file.File, defaults *defaults.Defaults, defaultConfig string) (tools.Tools, error) {
	if err := defaults.Validate(); err != nil {
		return nil, fmt.Errorf("validating defaults: %w", err)
	}

	loader := &ToolsLoader{
		File:     path,
		Defaults: defaults,
	}

	if err := loader.Load(defaultConfig); err != nil {
		return nil, fmt.Errorf("loading tools from %q: %w", path, err)
	}

	return loader.tools, nil
}

func (t *ToolsLoader) Length() int {
	return len(t.tools)
}

// LoadTools loads the tools configuration.
func (t *ToolsLoader) Load(defaultConfig string) error {
	data, err := t.Read()
	if err != nil {
		return fmt.Errorf("reading tools from %q: %w", t.File, err)
	}

	length, err := CountYamlListItems(data)
	if err != nil {
		return fmt.Errorf("counting items in tools from %q: %w", t.File, err)
	}

	inherits, err := GetInherits(data, defaultConfig, length)
	if err != nil {
		return fmt.Errorf("getting inherits from tools from %q: %w", t.File, err)
	}

	inherits = []inherit.Inherit(inherits)

	if t.tools, err = tools.NewToolsFromDefaults(t.Defaults, inherits); err != nil {
		return fmt.Errorf("creating tools (to set defaults) from %q: %w", t.File, err)
	}

	if err := t.UnmarshalPieceByPiece(data); err != nil {
		return fmt.Errorf("loading tools (to set defaults) from %q: %w", t.File, err)
	}

	return nil
}

// Type represents the type of tool configuration source.
type Type int

const (
	// File represents a file path.
	FILE Type = iota
	// Stdin represents standard input.
	STDIN
)

// Type returns the type of tool configuration source.
func (t *ToolsLoader) Type() (Type, error) {
	if t.File.Path() == "-" {
		if !stdin.IsPiped() {
			return STDIN, errors.New("no data piped to stdin")
		}

		return STDIN, nil
	}

	return FILE, nil
}

// Read reads a tool configuration file and loads it into the a Tools collection.
// If the path is "-", it reads from stdin.
// Else, it reads from the specified file path.
func (t *ToolsLoader) Read() (data []byte, err error) {
	switch tt, err := t.Type(); {
	case err != nil:
		return nil, err
	case tt == STDIN:
		input, err := stdin.Read()
		if err != nil {
			return nil, err
		}

		data = []byte(input)
	default:
		// Read the YAML configuration file from disk.
		input, err := t.File.Read()
		if err != nil {
			return nil, err
		}

		data = input
	}

	return data, nil
}

// UnmarshalPieceByPiece unmarshals the YAML list-item by list-item.
func (t *ToolsLoader) UnmarshalPieceByPiece(data []byte) (err error) {
	listOfAny := []any{}

	err = yaml.Unmarshal(data, &listOfAny)
	if err != nil {
		return err
	}

	for i, item := range listOfAny {
		// Convert the item back to YAML bytes
		byteItems, err := yaml.Marshal(item)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(byteItems, &t.tools[i])
		if err != nil {
			return err
		}
	}

	return nil
}

// UnmarshalWhole unmarshals the entire YAML configuration into the tools collection.
func (t *ToolsLoader) UnmarshalWhole(data []byte) (err error) {
	err = yaml.Unmarshal(data, &t.tools)
	if err != nil {
		return err
	}

	return nil
}

// CountYamlListItems counts the number of items in a YAML list.
func CountYamlListItems(data []byte) (int, error) {
	var items []any

	if err := yaml.Unmarshal(data, &items); err != nil {
		return 0, err
	}

	return len(items), nil
}

type Inherits struct {
	Inherit string
}

func (i *Inherits) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		return nil
	}

	type rawInherits Inherits

	return unmarshal.DecodeWithOptionalKnownFields(value, (*rawInherits)(i), false, structs.New(i).Name())
}

func GetInherits(data []byte, defaultInherit string, length int) ([]inherit.Inherit, error) {
	items := make([]Inherits, length)

	if err := yaml.Unmarshal(data, &items); err != nil {
		return nil, err
	}

	inherits := make([]inherit.Inherit, length)

	for i := range items {
		inherits[i] = inherit.Inherit(defaultInherit)
	}

	for i := range items {
		if val := items[i].Inherit; val != "" {
			inherits[i] = inherit.Inherit(val)
		}
	}

	return inherits, nil
}
