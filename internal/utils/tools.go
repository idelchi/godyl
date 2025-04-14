package utils

import (
	"errors"
	"fmt"

	"github.com/idelchi/go-next-tag/pkg/stdin"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/pkg/path/file"

	"gopkg.in/yaml.v3"
)

func LoadTools(path file.File, defaults tools.Defaults) (tools.Tools, error) {
	loader := &ToolsLoader{
		File:     path,
		Defaults: defaults,
	}

	if err := loader.Load(); err != nil {
		return nil, fmt.Errorf("loading tools from %q: %w", path, err)
	}

	return loader.tools, nil
}

type ToolsLoader struct {
	File file.File

	Defaults tools.Defaults
	tools    tools.Tools
}

func (t *ToolsLoader) Length() int {
	return len(t.tools)
}

// LoadTools loads the tools configuration.
func (t *ToolsLoader) Load() error {
	data, err := t.Read()
	if err != nil {
		return fmt.Errorf("reading tools from %q: %w", t.File, err)
	}

	if err := t.UnmarshalWhole(data); err != nil {
		return fmt.Errorf("unmarshaling tools (to get length) from %q: %w", t.File, err)
	}

	if t.tools, err = tools.NewTools(t.Defaults, t.Length()); err != nil {
		return fmt.Errorf("creating tools (to set defaults) from %q: %w", t.File, err)
	}

	if err := t.UnmarshalPieceByPiece(data); err != nil {
		return fmt.Errorf("loading tools (to set defaults) from %q: %w", t.File, err)
	}

	return nil
}

type ToolsType int

const (
	FILE ToolsType = iota
	STDIN
)

func (t *ToolsLoader) Type() (ToolsType, error) {
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

func (t *ToolsLoader) UnmarshalWhole(data []byte) (err error) {
	err = yaml.Unmarshal(data, &t.tools)
	if err != nil {
		return err
	}

	return nil
}
