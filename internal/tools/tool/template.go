// Package tool provides core functionality for managing tool configurations.
package tool

import (
	"fmt"
	"maps"

	"github.com/idelchi/godyl/internal/templates"
	"github.com/idelchi/godyl/internal/tools/sources"
)

// ToTemplateMap converts the Tool struct to a map suitable for templating.
// It adds any additional maps provided in the flatten argument to the template map.
func (t *Tool) ToTemplateMap(flatten ...map[string]any) map[string]any {
	templateMap := map[string]any{
		"Name":    t.Name,
		"Env":     t.Env,
		"Values":  t.Values,
		"Version": t.Version.Version,
		"Exe":     t.Exe.Name,
		"Output":  t.Output,
	}

	for _, o := range flatten {
		maps.Copy(templateMap, o)
	}

	return templateMap
}

// TemplateError wraps an error with additional context about the template name.
func TemplateError(err error, name string) error {
	return fmt.Errorf("applying template to %q: %w", name, err)
}

// TemplateFirst applies templating to various fields of the Tool struct, such as version, path, and checksum.
// It processes these fields using Go templates and updates them with the templated values.
func (t *Tool) TemplateFirst() error {
	values := t.ToTemplateMap(t.Platform.ToMap())

	if err := templates.ApplyAndSet(&t.Name, values); err != nil {
		return TemplateError(err, "name")
	}

	// Apply templating to Source.Type
	output, err := templates.Apply(t.Source.Type.String(), values)
	if err != nil {
		return TemplateError(err, "source.type")
	}

	t.Source.Type = sources.Type(output)

	// Apply templating to the Skip conditions
	for i := range t.Skip {
		err = templates.ApplyAndSet(&t.Skip[i].Condition, values)
		if err != nil {
			return TemplateError(err, "skip.condition")
		}
	}

	if err := templates.ApplyAndSet(&t.Output, values); err != nil {
		return TemplateError(err, "output")
	}

	return nil
}

// TemplateLast applies templating to the remaining fields of the Tool struct.
// Templates:
//   - Exe.Patterns
//   - Extensions
//   - Commands
//   - Hints patterns and weights
//   - URL
func (t *Tool) TemplateLast() error {
	values := t.ToTemplateMap(t.Platform.ToMap())

	// Apply templating to Exe.Patterns
	patterns := *t.Exe.Patterns
	for i := range patterns {
		if err := templates.ApplyAndSet(&patterns[i], values); err != nil {
			return err
		}
	}

	// Apply templating to Extensions
	extensions := *t.Extensions
	for i := range extensions {
		if err := templates.ApplyAndSet(&extensions[i], values); err != nil {
			return err
		}
	}

	// Apply templating to commands
	for i, cmd := range t.Commands.Commands {
		output, err := templates.Apply(cmd.String(), values)
		if err != nil {
			return err
		}

		t.Commands.Commands[i].From(output)
	}

	// Apply templating to Hints patterns and weights
	hints := *t.Hints
	for i := range hints {
		if err := templates.ApplyAndSet(&hints[i].Pattern, values); err != nil {
			return err
		}

		if err := templates.ApplyAndSet(&hints[i].Weight, values); err != nil {
			return err
		}
	}

	if err := templates.ApplyAndSet(&t.URL, values); err != nil {
		return err
	}

	return nil
}
