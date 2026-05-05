package skip_test

import (
	"testing"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/goccy/go-yaml"

	"github.com/idelchi/godyl/internal/tools/skip"
)

func makeCondition(t *testing.T, template string) skip.Condition {
	t.Helper()

	var c skip.Condition

	c.Condition.Set(template)

	return c
}

// ---------------------------------------------------------------------------
// Condition tests
// ---------------------------------------------------------------------------

func TestConditionTrue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		template string
		want     bool
	}{
		{
			name:     "template true evaluates to true",
			template: "true",
			want:     true,
		},
		{
			name:     "template false evaluates to false",
			template: "false",
			want:     false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c := makeCondition(t, tc.template)

			if err := c.Parse(); err != nil {
				t.Fatalf("Parse() unexpected error: %v", err)
			}

			got := c.True()
			if got != tc.want {
				t.Errorf("True() = %v, want %v (template %q)", got, tc.want, tc.template)
			}
		})
	}
}

func TestConditionUnmarshalYAML(t *testing.T) {
	t.Parallel()

	// Construct a Condition through YAML unmarshaling instead of makeCondition.
	input := heredoc.Doc(`
		condition: "true"
		reason: "yaml-driven"
	`)

	var c skip.Condition

	if err := yaml.Unmarshal([]byte(input), &c); err != nil {
		t.Fatalf("yaml.Unmarshal() unexpected error: %v", err)
	}

	if c.Reason != "yaml-driven" {
		t.Errorf("Reason = %q, want %q", c.Reason, "yaml-driven")
	}

	if err := c.Parse(); err != nil {
		t.Fatalf("Parse() unexpected error: %v", err)
	}

	if !c.True() {
		t.Errorf("True() = false, want true after parsing condition: true")
	}
}

// ---------------------------------------------------------------------------
// Skip.Has tests
// ---------------------------------------------------------------------------

func TestSkipHas(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		build func(t *testing.T) skip.Skip
		want  bool
	}{
		{
			name:  "nil Skip has no conditions",
			build: func(_ *testing.T) skip.Skip { return nil },
			want:  false,
		},
		{
			name:  "empty Skip has no conditions",
			build: func(_ *testing.T) skip.Skip { return skip.Skip{} },
			want:  false,
		},
		{
			name: "Skip with one condition has conditions",
			build: func(t *testing.T) skip.Skip { //nolint:thelper // not a test helper, it's a builder
				c := makeCondition(t, "true")

				c.Reason = "always"

				return skip.Skip{c}
			},
			want: true,
		},
		{
			name: "Skip with multiple conditions has conditions",
			build: func(t *testing.T) skip.Skip { //nolint:thelper // not a test helper, it's a builder
				c1 := makeCondition(t, "true")
				c2 := makeCondition(t, "false")

				return skip.Skip{c1, c2}
			},
			want: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s := tc.build(t)
			got := s.Has()

			if got != tc.want {
				t.Errorf("Has() = %v, want %v", got, tc.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Skip.Evaluate tests
// ---------------------------------------------------------------------------

func TestEvaluate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		conditions []string // template strings for each Condition
		reasons    []string // optional per-condition reason (empty string = no reason)
		wantLen    int      // expected length of the returned Skip
		wantHas    bool     // expected Has() result on returned Skip
		wantReason string   // if non-empty, assert got[0].Reason equals this
	}{
		{
			name:       "single true condition returns non-empty skip",
			conditions: []string{"true"},
			wantLen:    1,
			wantHas:    true,
		},
		{
			name:       "single false condition returns empty skip",
			conditions: []string{"false"},
			wantLen:    0,
			wantHas:    false,
		},
		{
			name:       "all false conditions return empty skip",
			conditions: []string{"false", "false"},
			wantLen:    0,
			wantHas:    false,
		},
		{
			name:       "all true conditions return full skip",
			conditions: []string{"true", "true"},
			wantLen:    2,
			wantHas:    true,
		},
		{
			name:       "mixed conditions return only true ones",
			conditions: []string{"true", "false", "true"},
			wantLen:    2,
			wantHas:    true,
		},
		{
			name:       "empty skip evaluates to empty",
			conditions: []string{},
			wantLen:    0,
			wantHas:    false,
		},
		{
			name:       "true condition preserves reason",
			conditions: []string{"true"},
			reasons:    []string{"skip for testing"},
			wantLen:    1,
			wantHas:    true,
			wantReason: "skip for testing",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s := make(skip.Skip, 0, len(tc.conditions))
			for i, tmpl := range tc.conditions {
				c := makeCondition(t, tmpl)

				if i < len(tc.reasons) {
					c.Reason = tc.reasons[i]
				}

				s = append(s, c)
			}

			got, err := s.Evaluate()
			if err != nil {
				t.Fatalf("Evaluate() unexpected error: %v", err)
			}

			if len(got) != tc.wantLen {
				t.Errorf("Evaluate() returned Skip of length %d, want %d", len(got), tc.wantLen)
			}

			if got.Has() != tc.wantHas {
				t.Errorf("Evaluate().Has() = %v, want %v", got.Has(), tc.wantHas)
			}

			if tc.wantReason != "" {
				if len(got) == 0 {
					t.Fatalf("Evaluate() returned empty Skip, cannot check Reason")
				}

				if got[0].Reason != tc.wantReason {
					t.Errorf("Evaluate()[0].Reason = %q, want %q", got[0].Reason, tc.wantReason)
				}

				if !got[0].True() {
					t.Errorf("Evaluate()[0].True() = false, want true for a condition that evaluated to true")
				}
			}
		})
	}
}

// TestEvaluateMultipleErrorInputs is a table-driven test that verifies
// Evaluate() returns an error for a variety of malformed condition templates.
// Each template must cause condition.Parse() to fail.
func TestEvaluateMultipleErrorInputs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		template string
	}{
		{
			// Unclosed YAML mapping is invalid YAML syntax.
			name:     "unclosed YAML mapping",
			template: "{bad: yaml: [",
		},
		{
			// A YAML mapping value cannot be decoded into a bool.
			name:     "YAML mapping cannot decode to bool",
			template: "{key: value}",
		},
		{
			// A YAML sequence cannot be decoded into a bool.
			name:     "YAML sequence cannot decode to bool",
			template: "[1, 2, 3]",
		},
		{
			// A plain string that is not a recognised boolean literal fails.
			name:     "plain string not a bool",
			template: "not-a-bool",
		},
		{
			// A quoted non-boolean string also fails.
			name:     "quoted non-boolean string",
			template: `"hello"`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c := makeCondition(t, tc.template)
			s := skip.Skip{c}

			if _, err := s.Evaluate(); err == nil {
				t.Errorf("Evaluate() with template %q expected error, got nil", tc.template)
			}
		})
	}
}

func TestEvaluateErrorPath(t *testing.T) {
	t.Parallel()

	// A template that is not valid YAML for a bool causes Parse() to fail.
	// yaml.Unmarshal of "{bad: yaml: [" into bool will error.
	c := makeCondition(t, "{bad: yaml: [")

	s := skip.Skip{c}

	if _, err := s.Evaluate(); err == nil {
		t.Fatal("Evaluate() expected error for unparseable condition template, got nil")
	}
}
