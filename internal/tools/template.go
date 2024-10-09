package tools

import (
	"bytes"
	"strconv"
	"text/template"

	"github.com/idelchi/godyl/internal/tools/sources"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// NormalizeValues ensures all keys in Values are upper-cased
func (t *Tool) NormalizeValues() {
	normalizedValues := make(map[string]any)
	c := cases.Title(language.English)

	// Iterate through the Values and convert keys to uppercase
	for key, value := range t.Values {
		upperKey := c.String(key)
		normalizedValues[upperKey] = value
	}

	// Replace the original Values map with the normalized one
	t.Values = normalizedValues
}

// ApplyTemplate applies Go templates to a string field using the Tool struct as data
func (t *Tool) ApplyTemplate(field string) (string, error) {
	var buf bytes.Buffer
	tmpl, err := template.New("tmpl").Parse(field)
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
	skip, err := t.ApplyTemplate(t.SkipTemplate)
	if err != nil {
		return err
	}
	t.Skip, err = strconv.ParseBool(skip)
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

	t.Exe, err = t.ApplyTemplate(t.Exe)
	if err != nil {
		return err
	}

	// Apply templating to Source.Commands (iterate over the command list)
	for i, cmd := range t.Source.Commands {
		output, err := t.ApplyTemplate(string(cmd))
		if err != nil {
			return err
		}
		t.Source.Commands[i] = sources.Command(output)
	}

	// Apply templating to Source.Commands (iterate over the command list)
	for i, hints := range t.Hints {
		output, err := t.ApplyTemplate(hints.Pattern)
		if err != nil {
			return err
		}
		t.Hints[i].Pattern = output

		output, err = t.ApplyTemplate(hints.WeightTemplate)
		if err != nil {
			return err
		}
		// Convert the result (string) into an integer and store it in the actual Weight field
		t.Hints[i].Weight, err = strconv.Atoi(output)
		if err != nil {
			return err
		}

	}

	return nil
}
