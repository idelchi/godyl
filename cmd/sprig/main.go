package main

import (
	"os"
	"text/template"

	sprig "github.com/go-task/slim-sprig/v3"
)

func main() {
	tmpl, _ := template.New("test").Funcs(sprig.FuncMap()).Parse(`
Platform OS: {{ .Platform.OS }}
Result: {{- if has .Platform.OS (list "linux" "darwin") -}}
    first
{{- else -}}
    second
{{- end -}}
`)

	data := struct {
		Platform struct {
			OS string
		}
	}{
		Platform: struct {
			OS string
		}{
			OS: "linux",
		},
	}

	err := tmpl.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}
