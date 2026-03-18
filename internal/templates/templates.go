package templates

import (
	"bytes"
	"text/template"

	sprig "github.com/go-task/slim-sprig/v3"
)

// Apply processes a template field with the provided values.
func Apply(field string, values any) (string, error) {
	var buf bytes.Buffer

	tmpl, err := template.New("tmpl").Funcs(sprig.FuncMap()).Option("missingkey=error").Parse(field)
	if err != nil {
		return "", err
	}

	if err := tmpl.Execute(&buf, values); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// ApplyAndSet processes a template field and updates the field in place.
// If Apply returns an error, the field is left unchanged.
func ApplyAndSet(field *string, values any) error {
	result, err := Apply(*field, values)
	if err == nil {
		*field = result
	}

	return err
}
