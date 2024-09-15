// Package templates provides template processing utilities for configuration and content generation.
package templates

import (
	"bytes"
	"maps"
	"text/template"

	sprig "github.com/go-task/slim-sprig/v3"
)

// Processor handles template processing with configurable options.
type Processor struct {
	funcMap template.FuncMap
	values  map[string]any
	option  string
}

// Option is a function that modifies a Processor.
type Option func(*Processor)

// WithMissingKeyDefault sets missing keys to display "<no value>".
func WithMissingKeyDefault() Option {
	return func(p *Processor) {
		p.option = "missingkey=default"
	}
}

// WithMissingKeyZero sets missing keys to use the zero value.
func WithMissingKeyZero() Option {
	return func(p *Processor) {
		p.option = "missingkey=zero"
	}
}

// WithMissingKeyError makes missing keys return an error.
func WithMissingKeyError() Option {
	return func(p *Processor) {
		p.option = "missingkey=error"
	}
}

// WithSlimSprig adds slim-sprig functions.
func WithSlimSprig() Option {
	return func(p *Processor) {
		p.funcMap = sprig.FuncMap()
	}
}

// New creates a processor with the given options.
func New(opts ...Option) *Processor {
	p := &Processor{
		option:  "missingkey=default", // default if not specified
		funcMap: template.FuncMap{},
		values:  make(map[string]any),
	}

	// Apply all options
	for _, opt := range opts {
		opt(p)
	}

	return p
}

// WithOptions adds options to the processor configuration.
func (p *Processor) WithOptions(opts ...Option) *Processor {
	for _, opt := range opts {
		opt(p)
	}

	return p
}

// WithValues adds value maps to the processor.
func (p *Processor) WithValues(values ...map[string]any) *Processor {
	p.AddValues(values...)

	return p
}

// AddValue adds or updates a single value.
func (p *Processor) AddValue(key string, value any) {
	p.values[key] = value
}

// Values returns the merged values map from the processor.
func (p *Processor) Values() map[string]any {
	return p.values
}

// AddValues adds or updates multiple values.
func (p *Processor) AddValues(values ...map[string]any) {
	for _, m := range values {
		maps.Copy(p.values, m)
	}
}

// AddFuncs adds or updates multiple functions.
func (p *Processor) AddFuncs(funcs ...template.FuncMap) {
	for _, m := range funcs {
		maps.Copy(p.funcMap, m)
	}
}

// Apply processes a template string with the accumulated values.
func (p *Processor) Apply(templateStr string) (string, error) {
	tmpl, err := template.New("tmpl").
		Funcs(p.funcMap).
		Option(p.option).
		Parse(templateStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, p.values); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// ApplyAndSet processes a template field and updates it in place.
func (p *Processor) ApplyAndSet(field *string) error {
	tmpl, err := p.Apply(*field)
	if err == nil {
		*field = tmpl
	}

	return err
}

// Reset clears all accumulated values.
func (p *Processor) Reset() {
	p.values = make(map[string]any)
}
