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

// YAML is a struct that provides methods to edit YAML files while preserving comments.
type YAML struct {
	// File is the YAML document being edited.
	File file.File

	// comments preserves comments across load and save operations.
	comments yaml.CommentMap

	// data contains the decoded YAML document.
	data map[string]any
}

// New creates a new YAML editor instance with the specified file.
func New(file file.File) *YAML {
	return &YAML{
		File: file,

		comments: yaml.CommentMap{},
	}
}

// Load decodes the YAML file and its comments, creating the file and parent directory when absent.
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

	if y.data == nil {
		y.data = make(map[string]any)
	}

	return nil
}

// Save writes the current data back to the YAML file, preserving comments.
func (y *YAML) Save() error {
	result, err := yaml.MarshalWithOptions(y.data, append(pretty.DefaultYAMLOptions(), yaml.WithComment(y.comments))...)
	if err != nil {
		return err
	}

	// Write back to file
	return y.File.Write(result)
}

// Write replaces the document data while preserving comments loaded from the existing file.
func (y *YAML) Write(input map[string]any) error {
	if err := y.Load(); err != nil {
		return err
	}

	y.data = input

	return y.Save()
}

// Merge overlays input on the existing YAML document while preserving its comments.
func (y *YAML) Merge(input map[string]any) error {
	if err := y.Load(); err != nil {
		return err
	}

	maps.Copy(y.data, input)

	return y.Save()
}
