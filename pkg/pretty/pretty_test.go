package pretty_test

import (
	"strings"
	"testing"

	"github.com/idelchi/godyl/pkg/pretty"
)

// maskedStruct has a field tagged with mask:"fixed".
// The go-mask default masker replaces the field value with 8 '*' characters.
type maskedStruct struct {
	Name   string `json:"name"   yaml:"name"`
	Secret string `json:"secret" mask:"fixed" yaml:"secret"`
}

// maskedValue is the fixed-width mask string produced by go-mask's "fixed" strategy.
// MaskJSON creates an instance Masker, calls SetMaskChar("-"), and registers both
// MaskFilledString ("filled") and MaskFixedString ("fixed") on that instance.
// MaskFixedString produces 8 repetitions of the masker's MaskChar, which is "-",
// so the result is 8 dashes.
const maskedValue = "--------"

func TestYAML(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    any
		want     string
		exact    bool
		contains string // additional substring that must appear
	}{
		{
			name:  "single-key map produces exact YAML output",
			input: map[string]string{"key": "value"},
			want:  "key: value\n",
			exact: true,
		},
		{
			name:  "nested struct fields appear in YAML output",
			input: maskedStruct{Name: "alice", Secret: "topsecret"},
			// Multi-field structs may vary in field order; use Contains for robustness.
			want:     "name: alice",
			contains: "topsecret",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := pretty.YAML(tc.input)

			if tc.exact {
				if got != tc.want {
					t.Errorf("YAML(%v) =\n%q\nwant\n%q", tc.input, got, tc.want)
				}
			} else {
				if !strings.Contains(got, tc.want) {
					t.Errorf("YAML(%v): expected output to contain %q, got:\n%s", tc.input, tc.want, got)
				}

				if tc.contains != "" && !strings.Contains(got, tc.contains) {
					t.Errorf("YAML(%v): expected output to contain raw value %q, got:\n%s", tc.input, tc.contains, got)
				}
			}
		})
	}
}

func TestJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    any
		want     string
		exact    bool
		contains string // additional substring that must appear
	}{
		{
			name:  "single-key map produces exact JSON output",
			input: map[string]string{"key": "value"},
			want:  "{\n      \"key\": \"value\"\n  }",
			exact: true,
		},
		{
			name:  "struct fields appear in JSON output",
			input: maskedStruct{Name: "bob", Secret: "password"},
			// Multi-field structs may vary; use Contains for robustness.
			want:     `"name": "bob"`,
			contains: "password",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := pretty.JSON(tc.input)

			if tc.exact {
				if got != tc.want {
					t.Errorf("JSON(%v) =\n%q\nwant\n%q", tc.input, got, tc.want)
				}
			} else {
				if !strings.Contains(got, tc.want) {
					t.Errorf("JSON(%v): expected output to contain %q, got:\n%s", tc.input, tc.want, got)
				}

				if tc.contains != "" && !strings.Contains(got, tc.contains) {
					t.Errorf("JSON(%v): expected output to contain raw value %q, got:\n%s", tc.input, tc.contains, got)
				}
			}
		})
	}
}

func TestYAMLMasked(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		input          maskedStruct
		notContains    string
		containsMasked string
		containsPlain  string
	}{
		{
			name:           "secret field is replaced with fixed mask",
			input:          maskedStruct{Name: "charlie", Secret: "supersecret"},
			notContains:    "supersecret",
			containsMasked: maskedValue,
			containsPlain:  "charlie",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := pretty.YAMLMasked(tc.input)

			if strings.Contains(got, tc.notContains) {
				t.Errorf("YAMLMasked: expected output NOT to contain %q, got:\n%s", tc.notContains, got)
			}

			if !strings.Contains(got, tc.containsMasked) {
				t.Errorf("YAMLMasked: expected output to contain masked value %q, got:\n%s", tc.containsMasked, got)
			}

			if !strings.Contains(got, tc.containsPlain) {
				t.Errorf("YAMLMasked: expected output to contain plain field %q, got:\n%s", tc.containsPlain, got)
			}
		})
	}
}

func TestJSONMasked(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		input          maskedStruct
		notContains    string
		containsMasked string
		containsPlain  string
	}{
		{
			name:           "secret field is replaced with fixed mask",
			input:          maskedStruct{Name: "dana", Secret: "mypassword"},
			notContains:    "mypassword",
			containsMasked: maskedValue,
			containsPlain:  "dana",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := pretty.JSONMasked(tc.input)

			if strings.Contains(got, tc.notContains) {
				t.Errorf("JSONMasked: expected output NOT to contain %q, got:\n%s", tc.notContains, got)
			}

			if !strings.Contains(got, tc.containsMasked) {
				t.Errorf("JSONMasked: expected output to contain masked value %q, got:\n%s", tc.containsMasked, got)
			}

			if !strings.Contains(got, tc.containsPlain) {
				t.Errorf("JSONMasked: expected output to contain plain field %q, got:\n%s", tc.containsPlain, got)
			}
		})
	}
}

func TestYAMLNilInput(t *testing.T) {
	t.Parallel()

	// YAML(nil) must not panic. The go-yaml library marshals nil as "null\n".
	// The function returns the error message as a string on failure, so any
	// non-empty string is an acceptable outcome — what matters is no panic.
	got := pretty.YAML(nil)

	if got == "" {
		t.Error("YAML(nil) returned empty string, want non-empty output (e.g. \"null\\n\")")
	}
}

func TestJSONNilInput(t *testing.T) {
	t.Parallel()

	// JSON(nil) must not panic. encoding/json marshals nil as "null".
	// The function returns the error message as a string on failure, so any
	// non-empty string is an acceptable outcome — what matters is no panic.
	got := pretty.JSON(nil)

	if got == "" {
		t.Error("JSON(nil) returned empty string, want non-empty output (e.g. \"null\")")
	}
}

func TestEnv(t *testing.T) {
	t.Parallel()

	// godotenv.Marshal quotes string values, so the format is KEY="value".
	tests := []struct {
		name     string
		input    any
		contains string
	}{
		{
			name:     "map entry appears in KEY=value dotenv format",
			input:    map[string]string{"KEY": "value"},
			contains: `KEY="value"`,
		},
		{
			name:     "lowercase key preserves case in dotenv output",
			input:    map[string]string{"mykey": "myval"},
			contains: `mykey="myval"`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := pretty.Env(tc.input)
			if !strings.Contains(got, tc.contains) {
				t.Errorf("Env(%v): expected output to contain %q, got:\n%s", tc.input, tc.contains, got)
			}
		})
	}
}
