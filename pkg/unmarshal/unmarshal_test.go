package unmarshal_test

import (
	"slices"
	"testing"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/goccy/go-yaml/parser"

	"github.com/idelchi/godyl/pkg/unmarshal"
)

// testStruct is used by TestStrict and TestLax.
type testStruct struct {
	Name  string `yaml:"name"`
	Value int    `yaml:"value"`
}

func TestStrict(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		yaml    string
		want    testStruct
		wantErr bool
	}{
		{
			name: "valid yaml",
			yaml: heredoc.Doc(`
				name: alice
				value: 42
			`),
			want:    testStruct{Name: "alice", Value: 42},
			wantErr: false,
		},
		{
			name: "unknown field",
			yaml: heredoc.Doc(`
				name: alice
				unknown: extra
			`),
			want:    testStruct{},
			wantErr: true,
		},
		{
			name:    "empty string",
			yaml:    "",
			want:    testStruct{},
			wantErr: false,
		},
		{
			name:    "partial yaml only name",
			yaml:    "name: bob",
			want:    testStruct{Name: "bob", Value: 0},
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var got testStruct

			err := unmarshal.Strict([]byte(tc.yaml), &got)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("Strict(%q): expected error, got nil", tc.yaml)
				}

				return
			}

			if err != nil {
				t.Fatalf("Strict(%q): unexpected error: %v", tc.yaml, err)
			}

			if got != tc.want {
				t.Errorf("Strict(%q): got %+v, want %+v", tc.yaml, got, tc.want)
			}
		})
	}
}

func TestLax(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		yaml    string
		want    testStruct
		wantErr bool
	}{
		{
			name: "valid yaml",
			yaml: heredoc.Doc(`
				name: alice
				value: 42
			`),
			want:    testStruct{Name: "alice", Value: 42},
			wantErr: false,
		},
		{
			name: "unknown field allowed",
			yaml: heredoc.Doc(`
				name: alice
				unknown: extra
			`),
			want:    testStruct{Name: "alice", Value: 0},
			wantErr: false,
		},
		{
			name:    "invalid yaml",
			yaml:    ":::bad",
			want:    testStruct{},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var got testStruct

			err := unmarshal.Lax([]byte(tc.yaml), &got)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("Lax(%q): expected error, got nil", tc.yaml)
				}

				return
			}

			if err != nil {
				t.Fatalf("Lax(%q): unexpected error: %v", tc.yaml, err)
			}

			if got != tc.want {
				t.Errorf("Lax(%q): got %+v, want %+v", tc.yaml, got, tc.want)
			}
		})
	}
}

func TestSingleOrSlice(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		yaml string
		want []string
	}{
		{
			name: "single string",
			yaml: "foo",
			want: []string{"foo"},
		},
		{
			name: "sequence",
			yaml: heredoc.Doc(`
				- foo
				- bar
			`),
			want: []string{"foo", "bar"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f, err := parser.ParseBytes([]byte(tc.yaml), 0)
			if err != nil {
				t.Fatalf("parser.ParseBytes(%q): %v", tc.yaml, err)
			}

			if len(f.Docs) == 0 || f.Docs[0].Body == nil {
				t.Fatalf("parser.ParseBytes(%q): no document body", tc.yaml)
			}

			node := f.Docs[0].Body

			got, err := unmarshal.SingleOrSlice[string](node)
			if err != nil {
				t.Fatalf("SingleOrSlice(%q): unexpected error: %v", tc.yaml, err)
			}

			if !slices.Equal(got, tc.want) {
				t.Errorf("SingleOrSlice(%q) = %v, want %v", tc.yaml, got, tc.want)
			}
		})
	}
}

// TestSingleOrSlice_TypeMismatch verifies that decoding a string sequence into
// a mismatched element type (int) returns an error. The type parameter difference
// is invisible in a shared table struct, so this case lives in its own test.
func TestSingleOrSlice_TypeMismatch(t *testing.T) {
	t.Parallel()

	yaml := heredoc.Doc(`
		- foo
		- bar
	`)

	f, err := parser.ParseBytes([]byte(yaml), 0)
	if err != nil {
		t.Fatalf("parser.ParseBytes(%q): %v", yaml, err)
	}

	if len(f.Docs) == 0 || f.Docs[0].Body == nil {
		t.Fatalf("parser.ParseBytes(%q): no document body", yaml)
	}

	// Decode into int — strings in the sequence cannot parse as int.
	_, err = unmarshal.SingleOrSlice[int](f.Docs[0].Body)
	if err == nil {
		t.Errorf("SingleOrSlice[int](%q): expected error, got nil", yaml)
	}
}

// named is used by TestSingleStringOrStruct.
// The Name field carries single:"true" so that a bare string shorthand populates it.
type named struct {
	Name  string `single:"true" yaml:"name"`
	Value string `yaml:"value"`
}

func TestSingleStringOrStruct(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		yaml string
		want named
	}{
		{
			name: "string shorthand",
			yaml: "myname",
			want: named{Name: "myname"},
		},
		{
			name: "full struct",
			yaml: heredoc.Doc(`
				name: myname
				value: myval
			`),
			want: named{Name: "myname", Value: "myval"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			f, err := parser.ParseBytes([]byte(tc.yaml), 0)
			if err != nil {
				t.Fatalf("parser.ParseBytes(%q): %v", tc.yaml, err)
			}

			if len(f.Docs) == 0 || f.Docs[0].Body == nil {
				t.Fatalf("parser.ParseBytes(%q): no document body", tc.yaml)
			}

			node := f.Docs[0].Body

			var out named
			if err := unmarshal.SingleStringOrStruct(node, &out); err != nil {
				t.Fatalf("SingleStringOrStruct(%q): unexpected error: %v", tc.yaml, err)
			}

			if out != tc.want {
				t.Errorf("SingleStringOrStruct(%q): got %+v, want %+v", tc.yaml, out, tc.want)
			}
		})
	}
}

// noSingleTag is a struct with no field carrying the single:"true" tag.
type noSingleTag struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// TestSingleStringOrStruct_NoSingleTag verifies that passing a scalar YAML node
// to SingleStringOrStruct when the target struct has no single:"true" field returns
// an error instead of silently doing nothing.
func TestSingleStringOrStruct_NoSingleTag(t *testing.T) {
	t.Parallel()

	const yaml = "myname"

	f, err := parser.ParseBytes([]byte(yaml), 0)
	if err != nil {
		t.Fatalf("parser.ParseBytes(%q): %v", yaml, err)
	}

	if len(f.Docs) == 0 || f.Docs[0].Body == nil {
		t.Fatalf("parser.ParseBytes(%q): no document body", yaml)
	}

	var out noSingleTag

	err = unmarshal.SingleStringOrStruct(f.Docs[0].Body, &out)
	if err == nil {
		t.Error("SingleStringOrStruct with no single:\"true\" tag: expected error, got nil")
	}
}

func TestStrictMalformed(t *testing.T) {
	t.Parallel()

	// ":::bad" is not valid YAML — the parser rejects it before any struct
	// field inspection takes place.
	var got testStruct

	err := unmarshal.Strict([]byte(":::bad"), &got)
	if err == nil {
		t.Error("Strict(\":::bad\"): expected error for malformed YAML, got nil")
	}
}

func TestSingleOrSliceEmptyNode(t *testing.T) {
	t.Parallel()

	// Parse "null" to obtain an *ast.NullNode — a non-nil node that represents
	// an absent value.  SingleOrSlice falls into the default branch, calls
	// Decode on the NullNode, and returns an empty (non-nil) slice with no error.
	f, err := parser.ParseBytes([]byte("null"), 0)
	if err != nil {
		t.Fatalf("parser.ParseBytes(\"null\"): %v", err)
	}

	if len(f.Docs) == 0 || f.Docs[0].Body == nil {
		t.Fatal("parser.ParseBytes(\"null\"): no document body")
	}

	got, err := unmarshal.SingleOrSlice[string](f.Docs[0].Body)
	if err != nil {
		t.Fatalf("SingleOrSlice[string](NullNode): unexpected error: %v", err)
	}

	// The NullNode decodes to a zero-value string "", so the result is [""].
	// Regardless of the exact element, what matters is no panic and no error.
	if got == nil {
		t.Error("SingleOrSlice[string](NullNode) = nil, want non-nil slice")
	}
}

func TestTemplatable(t *testing.T) {
	t.Parallel()

	t.Run("set and parse bool", func(t *testing.T) {
		t.Parallel()

		var tpl unmarshal.Templatable[bool]
		tpl.Set("true")

		if err := tpl.Parse(); err != nil {
			t.Fatalf("Parse(): unexpected error: %v", err)
		}

		got, err := tpl.Get()
		if err != nil {
			t.Fatalf("Get(): unexpected error: %v", err)
		}

		if !got {
			t.Errorf("Get(): got %v, want true", got)
		}
	})

	t.Run("unset", func(t *testing.T) {
		t.Parallel()

		var tpl unmarshal.Templatable[string]

		if !tpl.IsUnset() {
			t.Errorf("IsUnset(): got false for fresh Templatable, want true")
		}
	})

	t.Run("IsUnset after Set empty string", func(t *testing.T) {
		t.Parallel()

		var tpl unmarshal.Templatable[string]
		tpl.Set("")

		// Set("") assigns an empty Template, so IsUnset should still be true.
		if !tpl.IsUnset() {
			t.Errorf("IsUnset() after Set(\"\"): got false, want true")
		}
	})

	t.Run("get before parse returns error and zero value", func(t *testing.T) {
		t.Parallel()

		var tpl unmarshal.Templatable[string]
		tpl.Set("hello")

		// Parse has not been called, so Get must return an error.
		got, err := tpl.Get()
		if err == nil {
			t.Errorf("Get() before Parse(): expected error, got nil")
		}

		// The returned zero value for string is "".
		if got != "" {
			t.Errorf("Get() before Parse(): got %q, want zero value \"\"", got)
		}
	})

	t.Run("parse string", func(t *testing.T) {
		t.Parallel()

		var tpl unmarshal.Templatable[string]
		tpl.Set("hello")

		if err := tpl.Parse(); err != nil {
			t.Fatalf("Parse(): unexpected error: %v", err)
		}

		got, err := tpl.Get()
		if err != nil {
			t.Fatalf("Get(): unexpected error: %v", err)
		}

		if got != "hello" {
			t.Errorf("Get(): got %q, want %q", got, "hello")
		}
	})

	t.Run("parse invalid int", func(t *testing.T) {
		t.Parallel()

		var tpl unmarshal.Templatable[int]
		tpl.Set("not_a_number")

		if err := tpl.Parse(); err == nil {
			t.Errorf("Parse(%q) on Templatable[int]: expected error, got nil", "not_a_number")
		}
	})
}
