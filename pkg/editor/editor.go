// Package editor provides functionality to edit YAML files while preserving comments.
package editor

import (
	"maps"

	"github.com/goccy/go-yaml"

	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/folder"
	"github.com/idelchi/godyl/pkg/pretty"
	"github.com/idelchi/godyl/pkg/unmarshal"
)

// Editor is a struct that provides methods to edit YAML files while preserving comments.
type YAML struct {
	File file.File

	comments yaml.CommentMap

	data map[string]any
}

// New creates a new YAML editor instance with the specified file.
func New(file file.File) *YAML {
	return &YAML{
		File: file,

		comments: yaml.CommentMap{},
	}
}

// Load loads the YAML file along with its comments into the editor.
func (y *YAML) Load() error {
	if !y.File.Exists() {
		if err := folder.FromFile(y.File).Create(); err != nil {
			return err
		}

		if err := y.File.Create(); err != nil {
			return err
		}
	}

	// Read the file
	data, err := y.File.Read()
	if err != nil {
		return err
	}

	// Parse with comment preservation
	y.data = make(map[string]any)

	if err := unmarshal.Lax(data, &y.data, yaml.CommentToMap(y.comments)); err != nil {
		return err
	}

	return nil
}

// Save writes the current data back to the YAML file, preserving comments.
func (y *YAML) Save() error {
	result, err := yaml.MarshalWithOptions(y.data, append(pretty.YAMLOptions, yaml.WithComment(y.comments))...)
	if err != nil {
		return err
	}

	// Write back to file
	return y.File.Write(result)
}

// Write updates the YAML file with the provided map.
func (y *YAML) Write(input map[string]any) error {
	if err := y.Load(); err != nil {
		return err
	}

	y.data = input

	return y.Save()
}

// Merge updates the YAML file with the provided map.
// It reads the existing YAML file, and merges the input map into it.
func (y *YAML) Merge(input map[string]any) error {
	if err := y.Load(); err != nil {
		return err
	}

	maps.Copy(y.data, input)

	return y.Save()
}
