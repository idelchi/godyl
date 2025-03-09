// Package templates provides functionality for template processing and application.
package templates

import (
	"bytes"
	"fmt"
	"html/template"

	sprig "github.com/go-task/slim-sprig/v3"
)

// Apply processes a template string with the provided values and returns the result.
func Apply(field string, values any) (string, error) {
	var buf bytes.Buffer

	tmpl, err := template.New("tmpl").Funcs(sprig.FuncMap()).Parse(field)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	if err := tmpl.Execute(&buf, values); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// ApplyAndSet processes a template string with the provided values and updates the field pointer.
func ApplyAndSet(field *string, values any) (err error) {
	*field, err = Apply(*field, values)

	return
}
