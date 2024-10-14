package tools

import (
	"bytes"
	"strconv"
	"text/template"

	sprig "github.com/go-task/slim-sprig/v3"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/pkg/utils"
)

// NormalizeValues ensures all keys in Values are capitalized.
func (t *Tool) NormalizeValues() {
	t.Values = utils.NormalizeMap(t.Values)
}

// ApplyTemplate applies Go templates to a string field using the Tool struct as data
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

// Template applies templating to the Tool's fields
func (t *Tool) Template() error {
	var err error

	// Apply templating to all relevant fields
	t.Name, err = t.ApplyTemplate(t.Name)
	if err != nil {
		return err
	}

	// Apply templating to all relevant fields
	t.Source.Type, err = t.ApplyTemplate(t.Source.Type)
	if err != nil {
		return err
	}

	for i, pattern := range t.Skip {
		t.Skip[i].Condition, err = t.ApplyTemplate(pattern.Condition)
		if err != nil {
			return err
		}
	}

	_, _, err = t.Skip.IsSkipped()
	if err != nil {
		return err
	}

	t.Version, err = t.ApplyTemplate(t.Version)
	if err != nil {
		return err
	}

	t.Path, err = t.ApplyTemplate(t.Path)
	if err != nil {
		return err
	}

	t.Checksum, err = t.ApplyTemplate(t.Checksum)
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

	for i, pattern := range t.Exe.Patterns {
		t.Exe.Patterns[i], err = t.ApplyTemplate(pattern)
		if err != nil {
			return err
		}
	}

	for i, cmd := range t.Source.Commands {
		output, err := t.ApplyTemplate(string(cmd))
		if err != nil {
			return err
		}
		t.Source.Commands[i] = sources.Command(output)
	}

	for i, cmd := range t.Post {
		output, err := t.ApplyTemplate(string(cmd))
		if err != nil {
			return err
		}
		t.Post[i] = sources.Command(output)
	}

	for i, hints := range t.Hints {
		t.Hints[i].Pattern, err = t.ApplyTemplate(hints.Pattern)
		if err != nil {
			return err
		}
		output, err := t.ApplyTemplate(hints.WeightTemplate)
		if err != nil {
			return err
		}
		// Convert the result (string) into an integer and store it in the actual Weight field
		utils.SetIfEmpty(&output, "1")
		t.Hints[i].Weight, err = strconv.Atoi(output)
		if err != nil {
			return err
		}
	}

	return nil
}
