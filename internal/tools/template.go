package tools

import (
	"bytes"
	"text/template"

	sprig "github.com/go-task/slim-sprig/v3"
	"github.com/idelchi/godyl/pkg/utils"
)

// NormalizeValues ensures that all keys in the Values map are capitalized.
// This helps to standardize the keys for consistent access and usage.
func (t *Tool) NormalizeValues() {
	t.Values = utils.NormalizeMap(t.Values)
}

// ApplyTemplate applies Go templates to a given string field using the Tool struct as the template data.
// It uses the slim-sprig library to provide additional template functions.
func (t *Tool) ApplyTemplate(field string) (string, error) {
	var buf bytes.Buffer
	tmpl, err := template.New("tmpl").Funcs(sprig.FuncMap()).Parse(field)
	if err != nil {
		return "", err
	}
	if err := tmpl.Execute(&buf, t); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Template applies templating to various fields of the Tool struct, such as version, path, and checksum.
// It processes these fields using Go templates and updates them with the templated values.
func (t *Tool) Template() error {
	var err error

	// Apply templating to Source.Type
	output, err := t.ApplyTemplate(t.Source.Type.String())
	if err != nil {
		return err
	}
	t.Source.Type.From(output)

	// Apply templating to the Skip conditions
	for i, pattern := range t.Skip {
		t.Skip[i].Condition, err = t.ApplyTemplate(pattern.Condition)
		if err != nil {
			return err
		}
	}

	// Validate the Skip conditions
	_, _, err = t.Skip.True()
	if err != nil {
		return err
	}

	// Apply templating to various fields like Version, Path, Checksum, Output, and Exe.Name
	t.Version, err = t.ApplyTemplate(t.Version)
	if err != nil {
		return err
	}

	t.Path, err = t.ApplyTemplate(t.Path)
	if err != nil {
		return err
	}

	t.Output, err = t.ApplyTemplate(t.Output)
	if err != nil {
		return err
	}

	t.Exe.Name, err = t.ApplyTemplate(t.Exe.Name)
	if err != nil {
		return err
	}

	// Apply templating to Exe.Patterns
	for i, pattern := range t.Exe.Patterns {
		t.Exe.Patterns[i], err = t.ApplyTemplate(pattern)
		if err != nil {
			return err
		}
	}

	// Apply templating to Source.Commands
	for i, cmd := range t.Source.Commands {
		output, err := t.ApplyTemplate(cmd.String())
		if err != nil {
			return err
		}
		t.Source.Commands[i].From(output)
	}

	// Apply templating to Post commands
	for i, cmd := range t.Post {
		output, err := t.ApplyTemplate(cmd.String())
		if err != nil {
			return err
		}
		t.Post[i].From(output)
	}

	// Apply templating to Hints patterns and weights
	for i, hints := range t.Hints {
		t.Hints[i].Pattern, err = t.ApplyTemplate(hints.Pattern)
		if err != nil {
			return err
		}
		t.Hints[i].Weight, err = t.ApplyTemplate(hints.Weight)
		if err != nil {
			return err
		}
		// Set a default weight of "1" if not specified and convert it to an integer
		utils.SetIfEmpty(&t.Hints[i].Weight, "1")

		if err := t.Hints[i].SetWeight(); err != nil {
			return err
		}
	}

	return nil
}
