// Package tool provides core functionality for managing tool configurations.
package tool

import (
	"fmt"

	"github.com/idelchi/godyl/internal/templates"
	"github.com/idelchi/godyl/internal/tools/sources"
)

// ToTemplateMap converts the Tool struct to a map suitable for templating.
// It adds any additional maps provided in the flatten argument to the template map.
func (t *Tool) ToTemplateMap(_ ...map[string]any) map[string]any {
	templateMap := map[string]any{
		"Name":   t.Name,
		"Env":    t.Env,
		"Values": t.Values,
		"Exe":    t.Exe.Name,
		"Source": t.Source.Type.String(),
		"Output": t.Output,
		"Tokens": map[string]string{
			"GitHub": t.Source.GitHub.Token,
			"GitLab": t.Source.GitLab.Token,
			"URL":    t.Source.URL.Token,
		},
	}

	if t.Version.Version != "" {
		templateMap["Version"] = t.Version.Version
	}

	return templateMap
}

// TemplateError wraps an error with additional context about the template name.
func TemplateError(err error, name string) error {
	return fmt.Errorf("applying template to %q: %w", name, err)
}

// TemplateFirst applies templating to various fields of the Tool struct, such as version, path, and checksum.
// It processes these fields using Go templates and updates them with the templated values.
func (t *Tool) TemplateFirst(tmpl *templates.Processor) error {
	if err := tmpl.ApplyAndSet(&t.Name); err != nil {
		return TemplateError(err, "name")
	}

	// Apply templating to Source.Type
	output, err := tmpl.Apply(t.Source.Type.String())
	if err != nil {
		return TemplateError(err, "source.type")
	}

	t.Source.Type = sources.Type(output)

	// Apply templating to the Skip conditions
	for i := range t.Skip {
		err = tmpl.ApplyAndSet(&t.Skip[i].Condition.Template)
		if err != nil {
			return TemplateError(err, "skip.condition")
		}
	}

	if err := tmpl.ApplyAndSet(&t.Output); err != nil {
		return TemplateError(err, "output")
	}

	if err := tmpl.ApplyAndSet(&t.Source.GitHub.Token); err != nil {
		return TemplateError(err, "github.token")
	}

	if err := tmpl.ApplyAndSet(&t.Source.GitLab.Token); err != nil {
		return TemplateError(err, "gitlab.token")
	}

	if err := tmpl.ApplyAndSet(&t.Source.URL.Token); err != nil {
		return TemplateError(err, "url.token")
	}

	return nil
}

// TemplateLast applies templating to the remaining fields of the Tool struct.
// Templates:
//   - Exe.Patterns
//   - Commands
//   - Hints patterns and weights
//   - URL
func (t *Tool) TemplateLast(tmpl *templates.Processor) error {
	// Apply templating to the url headers
	for key, value := range t.Source.URL.Headers {
		for i := range value {
			err := tmpl.ApplyAndSet(&value[i])
			if err != nil {
				return TemplateError(err, "url.headers."+key)
			}
		}
	}

	// Apply templating to Exe.Patterns
	patterns := *t.Exe.Patterns
	for i := range patterns {
		if err := tmpl.ApplyAndSet(&patterns[i]); err != nil {
			return err
		}
	}

	// Apply templating to commands
	for i, cmd := range t.Commands.Commands {
		output, err := tmpl.Apply(cmd.String())
		if err != nil {
			return err
		}

		t.Commands.Commands[i].From(output)
	}

	// Apply templating to Hints patterns and weights
	hints := *t.Hints
	for i := range hints {
		if err := tmpl.ApplyAndSet(&hints[i].Pattern); err != nil {
			return err
		}

		if err := tmpl.ApplyAndSet(&hints[i].Weight.Template); err != nil {
			return err
		}

		if err := tmpl.ApplyAndSet(&hints[i].Match.Template); err != nil {
			return err
		}
	}

	if err := tmpl.ApplyAndSet(&t.URL); err != nil {
		return err
	}

	return nil
}
