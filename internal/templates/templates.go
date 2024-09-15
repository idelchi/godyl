package templates

import (
	"bytes"
	"html/template"

	sprig "github.com/go-task/slim-sprig/v3"
)

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

func ApplyAndSet(field *string, values any) (err error) {
	*field, err = Apply(*field, values)

	return err
}
