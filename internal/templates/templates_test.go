package templates_test

import (
	"testing"

	"github.com/idelchi/godyl/internal/templates"
)

func TestApply(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		field     string
		values    any
		want      string
		wantError bool
	}{
		{
			name:   "simple substitution",
			field:  "{{ .Name }}",
			values: map[string]any{"Name": "foo"},
			want:   "foo",
		},
		{
			name:      "missing key errors",
			field:     "{{ .Missing }}",
			values:    map[string]any{},
			wantError: true,
		},
		{
			name:   "conditional true",
			field:  "{{ if .X }}yes{{ else }}no{{ end }}",
			values: map[string]any{"X": true},
			want:   "yes",
		},
		{
			name:   "conditional false",
			field:  "{{ if .X }}yes{{ else }}no{{ end }}",
			values: map[string]any{"X": false},
			want:   "no",
		},
		{
			name:   "nested access",
			field:  "{{ .A.B }}",
			values: map[string]any{"A": map[string]any{"B": "deep"}},
			want:   "deep",
		},
		{
			name:   "empty template",
			field:  "",
			values: nil,
			want:   "",
		},
		{
			name:   "literal text",
			field:  "no templates here",
			values: nil,
			want:   "no templates here",
		},
		{
			name:      "invalid syntax",
			field:     "{{ .Broken",
			values:    nil,
			wantError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := templates.Apply(tc.field, tc.values)

			if tc.wantError {
				if err == nil {
					t.Errorf("Apply(%q, %v): expected error, got nil", tc.field, tc.values)
				}

				return
			}

			if err != nil {
				t.Fatalf("Apply(%q, %v): unexpected error: %v", tc.field, tc.values, err)
			}

			if got != tc.want {
				t.Errorf("Apply(%q, %v): got %q, want %q", tc.field, tc.values, got, tc.want)
			}
		})
	}
}

func TestApplyAndSet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		field     string
		values    any
		want      string
		wantError bool
	}{
		{
			name:   "updates field in place",
			field:  "{{ .Name }}",
			values: map[string]any{"Name": "replaced"},
			want:   "replaced",
		},
		{
			name:      "error leaves field unchanged",
			field:     "{{ .Missing }}",
			values:    map[string]any{},
			wantError: true,
		},
		{
			name:   "literal text is preserved",
			field:  "static",
			values: nil,
			want:   "static",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			field := tc.field
			err := templates.ApplyAndSet(&field, tc.values)

			if tc.wantError {
				if err == nil {
					t.Errorf("ApplyAndSet: expected error, got nil")
				}

				if field != tc.field {
					t.Errorf("ApplyAndSet: field mutated on error: got %q, want original %q", field, tc.field)
				}

				return
			}

			if err != nil {
				t.Fatalf("ApplyAndSet: unexpected error: %v", err)
			}

			if field != tc.want {
				t.Errorf("ApplyAndSet: field = %q, want %q", field, tc.want)
			}
		})
	}
}

func TestProcessorApply(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		opts        []templates.Option
		addValues   map[string]any
		templateStr string
		want        string
		wantError   bool
	}{
		{
			name:        "missing key default shows no value",
			opts:        []templates.Option{templates.WithMissingKeyDefault()},
			templateStr: "{{ .Missing }}",
			want:        "<no value>",
		},
		{
			name:        "missing key zero renders nil as no value",
			opts:        []templates.Option{templates.WithMissingKeyZero()},
			templateStr: "{{ .Missing }}",
			want:        "<no value>",
		},
		{
			name:        "missing key error returns error",
			opts:        []templates.Option{templates.WithMissingKeyError()},
			templateStr: "{{ .Missing }}",
			wantError:   true,
		},
		{
			name:        "processor applies added values",
			opts:        []templates.Option{templates.WithMissingKeyError()},
			addValues:   map[string]any{"Name": "bar"},
			templateStr: "{{ .Name }}",
			want:        "bar",
		},
		{
			name:        "default option is missingkey=default",
			templateStr: "{{ .Missing }}",
			want:        "<no value>",
		},
		{
			name:        "empty template returns empty string",
			templateStr: "",
			want:        "",
		},
		{
			name:        "literal text unchanged",
			templateStr: "hello world",
			want:        "hello world",
		},
		{
			name:        "invalid syntax returns error",
			templateStr: "{{ .Broken",
			wantError:   true,
		},
		{
			name:        "sprig trimPrefix strips prefix",
			opts:        []templates.Option{templates.WithSlimSprig()},
			templateStr: `{{ "v1.2.3" | trimPrefix "v" }}`,
			want:        "1.2.3",
		},
		{
			name:        "sprig upper converts to uppercase",
			opts:        []templates.Option{templates.WithSlimSprig()},
			templateStr: `{{ "hello" | upper }}`,
			want:        "HELLO",
		},
		{
			name:        "sprig lower converts to lowercase",
			opts:        []templates.Option{templates.WithSlimSprig()},
			templateStr: `{{ "WORLD" | lower }}`,
			want:        "world",
		},
		{
			name:        "sprig trim removes whitespace",
			opts:        []templates.Option{templates.WithSlimSprig()},
			templateStr: `{{ "  spaces  " | trim }}`,
			want:        "spaces",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			p := templates.New(tc.opts...)
			if tc.addValues != nil {
				p.AddValues(tc.addValues)
			}

			got, err := p.Apply(tc.templateStr)

			if tc.wantError {
				if err == nil {
					t.Errorf("Processor.Apply(%q): expected error, got nil", tc.templateStr)
				}

				return
			}

			if err != nil {
				t.Fatalf("Processor.Apply(%q): unexpected error: %v", tc.templateStr, err)
			}

			if got != tc.want {
				t.Errorf("Processor.Apply(%q): got %q, want %q", tc.templateStr, got, tc.want)
			}
		})
	}
}

func TestProcessorApplyAndSet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		opts      []templates.Option
		values    map[string]any
		field     string
		want      string
		wantError bool
	}{
		{
			name:   "updates field with resolved value",
			opts:   []templates.Option{templates.WithMissingKeyError()},
			values: map[string]any{"Greeting": "hello"},
			field:  "{{ .Greeting }}",
			want:   "hello",
		},
		{
			name:      "error on missing key does not update field",
			opts:      []templates.Option{templates.WithMissingKeyError()},
			field:     "{{ .Missing }}",
			wantError: true,
		},
		{
			name:   "literal field is preserved",
			values: map[string]any{},
			field:  "unchanged",
			want:   "unchanged",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			p := templates.New(tc.opts...)
			if tc.values != nil {
				p.AddValues(tc.values)
			}

			field := tc.field
			err := p.ApplyAndSet(&field)

			if tc.wantError {
				if err == nil {
					t.Errorf("Processor.ApplyAndSet: expected error, got nil")
				}

				if field != tc.field {
					t.Errorf("Processor.ApplyAndSet: field mutated on error: got %q, want original %q", field, tc.field)
				}

				return
			}

			if err != nil {
				t.Fatalf("Processor.ApplyAndSet: unexpected error: %v", err)
			}

			if field != tc.want {
				t.Errorf("Processor.ApplyAndSet: field = %q, want %q", field, tc.want)
			}
		})
	}
}

func TestProcessorAddValues(t *testing.T) {
	t.Parallel()

	t.Run("AddValue makes key accessible", func(t *testing.T) {
		t.Parallel()

		p := templates.New(templates.WithMissingKeyError())
		p.AddValue("Key", "value1")

		got, err := p.Apply("{{ .Key }}")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != "value1" {
			t.Errorf("got %q, want %q", got, "value1")
		}
	})

	t.Run("AddValues merges all keys", func(t *testing.T) {
		t.Parallel()

		p := templates.New(templates.WithMissingKeyError())
		p.AddValues(
			map[string]any{"First": "one"},
			map[string]any{"Second": "two"},
		)

		got, err := p.Apply("{{ .First }}-{{ .Second }}")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != "one-two" {
			t.Errorf("got %q, want %q", got, "one-two")
		}
	})

	t.Run("AddValue overwrites existing key", func(t *testing.T) {
		t.Parallel()

		p := templates.New(templates.WithMissingKeyError())
		p.AddValue("X", "first")
		p.AddValue("X", "second")

		got, err := p.Apply("{{ .X }}")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != "second" {
			t.Errorf("got %q, want %q", got, "second")
		}
	})

	t.Run("AddValues overwrites existing key from prior map", func(t *testing.T) {
		t.Parallel()

		p := templates.New(templates.WithMissingKeyError())
		p.AddValues(
			map[string]any{"Key": "original"},
			map[string]any{"Key": "overwritten"},
		)

		got, err := p.Apply("{{ .Key }}")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != "overwritten" {
			t.Errorf("got %q, want %q", got, "overwritten")
		}
	})

	t.Run("WithValues returns processor for chaining", func(t *testing.T) {
		t.Parallel()

		p := templates.New(templates.WithMissingKeyError()).
			WithValues(map[string]any{"Chain": "chained"})

		got, err := p.Apply("{{ .Chain }}")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != "chained" {
			t.Errorf("got %q, want %q", got, "chained")
		}
	})
}

func TestApply_SprigAvailable(t *testing.T) {
	t.Parallel()

	// The package-level Apply uses sprig functions; verify trimPrefix is available.
	got, err := templates.Apply(`{{ trimPrefix "v" "v1.0" }}`, nil)
	if err != nil {
		t.Fatalf("Apply() unexpected error: %v", err)
	}

	if got != "1.0" {
		t.Errorf("Apply() = %q, want %q", got, "1.0")
	}
}

func TestProcessorWithOptions(t *testing.T) {
	t.Parallel()

	// Start with the default (missingkey=default), then switch to missingkey=error.
	p := templates.New()

	// Default option: missing key renders as "<no value>".
	got, err := p.Apply("{{ .Missing }}")
	if err != nil {
		t.Fatalf("Apply() with default option unexpected error: %v", err)
	}

	if got != "<no value>" {
		t.Errorf("Apply() with default option = %q, want %q", got, "<no value>")
	}

	// Switch to missingkey=error via WithOptions.
	p.WithOptions(templates.WithMissingKeyError())

	_, err = p.Apply("{{ .Missing }}")
	if err == nil {
		t.Error("Apply() after WithOptions(WithMissingKeyError()): expected error, got nil")
	}
}

func TestProcessorReset(t *testing.T) {
	t.Parallel()

	t.Run("reset clears all values", func(t *testing.T) {
		t.Parallel()

		p := templates.New()
		p.AddValue("Name", "before")
		p.AddValue("Other", "data")

		p.Reset()

		vals := p.Values()
		if len(vals) != 0 {
			t.Errorf("after Reset, Values() has %d entries, want 0", len(vals))
		}
	})

	t.Run("apply after reset uses empty values", func(t *testing.T) {
		t.Parallel()

		p := templates.New(templates.WithMissingKeyDefault())
		p.AddValue("Name", "present")
		p.Reset()

		got, err := p.Apply("{{ .Name }}")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != "<no value>" {
			t.Errorf("got %q, want %q", got, "<no value>")
		}
	})

	t.Run("values can be re-added after reset", func(t *testing.T) {
		t.Parallel()

		p := templates.New(templates.WithMissingKeyError())
		p.AddValue("X", "old")
		p.Reset()
		p.AddValue("X", "new")

		got, err := p.Apply("{{ .X }}")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != "new" {
			t.Errorf("got %q, want %q", got, "new")
		}
	})
}
