package generic_test

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/idelchi/godyl/pkg/generic"
)

func TestAnyNil(t *testing.T) {
	t.Parallel()

	one := 1
	two := 2

	tests := []struct {
		name string
		ptrs []*int
		want bool
	}{
		{
			name: "one nil",
			ptrs: []*int{&one, nil, &two},
			want: true,
		},
		{
			name: "none nil",
			ptrs: []*int{&one, &two},
			want: false,
		},
		{
			name: "all nil",
			ptrs: []*int{nil, nil},
			want: true,
		},
		{
			name: "empty",
			ptrs: []*int{},
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := generic.AnyNil(tc.ptrs...)
			if got != tc.want {
				t.Errorf("AnyNil(%v) = %v, want %v", tc.ptrs, got, tc.want)
			}
		})
	}
}

func TestIsZero(t *testing.T) {
	t.Parallel()

	type point struct{ X, Y int }

	tests := []struct {
		name string
		fn   func() bool
		want bool
	}{
		{name: "int zero", fn: func() bool { return generic.IsZero(0) }, want: true},
		{name: "int non-zero", fn: func() bool { return generic.IsZero(1) }, want: false},
		{name: "string empty", fn: func() bool { return generic.IsZero("") }, want: true},
		{name: "string non-empty", fn: func() bool { return generic.IsZero("x") }, want: false},
		{name: "bool false", fn: func() bool { return generic.IsZero(false) }, want: true},
		{name: "bool true", fn: func() bool { return generic.IsZero(true) }, want: false},
		{name: "struct zero value", fn: func() bool { return generic.IsZero(point{}) }, want: true},
		{name: "struct non-zero value", fn: func() bool { return generic.IsZero(point{X: 1}) }, want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.fn(); got != tc.want {
				t.Errorf("IsZero() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestSetIfZero(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		initial int
		value   int
		want    int
	}{
		{
			name:    "zero value gets set",
			initial: 0,
			value:   42,
			want:    42,
		},
		{
			name:    "non-zero value stays unchanged",
			initial: 10,
			value:   42,
			want:    10,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			v := tc.initial
			generic.SetIfZero(&v, tc.value)

			if v != tc.want {
				t.Errorf("after SetIfZero(%v, %v): got %v, want %v", tc.initial, tc.value, v, tc.want)
			}
		})
	}
}

func TestExpandHome(t *testing.T) {
	t.Parallel()

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("os.UserHomeDir(): %v", err)
	}

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "tilde alone",
			input: "~",
			want:  home,
		},
		{
			name:  "tilde with path",
			input: "~/some/path",
			want:  filepath.Join(home, "some/path"),
		},
		{
			name:  "no tilde",
			input: "/absolute/path",
			want:  "/absolute/path",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			// "~foo" starts with ~ but is not "~" and not "~/…".
			// ExpandHome treats path[1:] = "foo" as a relative segment and
			// joins it with the home directory: filepath.Join(home, "foo").
			// It does NOT expand to a different user's home; that is not supported.
			name:  "tilde not followed by slash treated as relative to home",
			input: "~foo",
			want:  filepath.Join(home, "foo"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := generic.ExpandHome(tc.input)
			if got != tc.want {
				t.Errorf("ExpandHome(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestDeepCopy(t *testing.T) {
	t.Parallel()

	type nested struct {
		Values []int
	}

	type payload struct {
		Name   string
		Nested nested
	}

	original := payload{
		Name:   "original",
		Nested: nested{Values: []int{1, 2, 3}},
	}

	copied, err := generic.DeepCopy(original)
	if err != nil {
		t.Fatalf("DeepCopy(): unexpected error: %v", err)
	}

	// Mutate the copy and verify the copy took effect.
	copied.Name = "mutated"
	copied.Nested.Values[0] = 99

	if copied.Name != "mutated" {
		t.Errorf("DeepCopy copy.Name: got %q, want \"mutated\"", copied.Name)
	}

	// Original must be unchanged.
	if original.Name != "original" {
		t.Errorf("DeepCopy mutated original.Name: got %q, want \"original\"", original.Name)
	}

	if original.Nested.Values[0] != 1 {
		t.Errorf("DeepCopy mutated original.Nested.Values[0]: got %d, want 1", original.Nested.Values[0])
	}
}

func TestDeepCopyPtr(t *testing.T) {
	t.Parallel()

	t.Run("nil input returns nil", func(t *testing.T) {
		t.Parallel()

		got, err := generic.DeepCopyPtr[int](nil)
		if err != nil {
			t.Fatalf("DeepCopyPtr(nil): unexpected error: %v", err)
		}

		if got != nil {
			t.Errorf("DeepCopyPtr(nil): got %v, want nil", got)
		}
	})

	t.Run("non-nil input returns distinct pointer with equal content", func(t *testing.T) {
		t.Parallel()

		v := 42
		src := &v

		dst, err := generic.DeepCopyPtr(src)
		if err != nil {
			t.Fatalf("DeepCopyPtr(&42): unexpected error: %v", err)
		}

		if dst == nil {
			t.Fatal("DeepCopyPtr(&42): got nil, want non-nil")
		}

		if dst == src {
			t.Error("DeepCopyPtr(&42): returned same pointer, want distinct pointer")
		}

		if *dst != *src {
			t.Errorf("DeepCopyPtr(&42): *dst = %d, want %d", *dst, *src)
		}
	})
}

func TestIsURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "https URL", input: "https://example.com", want: true},
		{name: "http URL", input: "http://example.com", want: true},
		{name: "ftp URL", input: "ftp://files.example.com", want: true},
		{name: "no scheme", input: "example.com", want: false},
		{name: "empty string", input: "", want: false},
		{name: "absolute path", input: "/foo/bar", want: false},
		{name: "file scheme no host", input: "file:///local/path", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := generic.IsURL(tc.input)
			if got != tc.want {
				t.Errorf("IsURL(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestIsSliceNilOrEmpty(t *testing.T) {
	t.Parallel()

	empty := []int{}
	nonEmpty := []int{1, 2, 3}

	tests := []struct {
		name string
		ptr  *[]int
		want bool
	}{
		{name: "nil pointer", ptr: nil, want: true},
		{name: "empty slice", ptr: &empty, want: true},
		{name: "non-empty slice", ptr: &nonEmpty, want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := generic.IsSliceNilOrEmpty(tc.ptr)
			if got != tc.want {
				t.Errorf("IsSliceNilOrEmpty(%v) = %v, want %v", tc.ptr, got, tc.want)
			}
		})
	}
}

func TestSafeDereference(t *testing.T) {
	t.Parallel()

	val := 99

	tests := []struct {
		name string
		ptr  *int
		want int
	}{
		{name: "non-nil pointer", ptr: &val, want: 99},
		{name: "nil pointer returns zero value", ptr: nil, want: 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := generic.SafeDereference(tc.ptr)
			if got != tc.want {
				t.Errorf("SafeDereference(%v) = %v, want %v", tc.ptr, got, tc.want)
			}
		})
	}
}

func TestPickByIndices(t *testing.T) {
	t.Parallel()

	source := []string{"a", "b", "c", "d"}

	tests := []struct {
		name    string
		indices []int
		want    []string
	}{
		{
			name:    "first and last",
			indices: []int{0, 3},
			want:    []string{"a", "d"},
		},
		{
			name:    "middle two",
			indices: []int{1, 2},
			want:    []string{"b", "c"},
		},
		{
			name:    "empty indices",
			indices: []int{},
			want:    []string{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := generic.PickByIndices(source, tc.indices)

			if !slices.Equal(got, tc.want) {
				t.Errorf("PickByIndices(%v, %v) = %v, want %v", source, tc.indices, got, tc.want)
			}
		})
	}

	t.Run("out-of-bounds panics", func(t *testing.T) {
		t.Parallel()

		defer func() {
			if r := recover(); r == nil {
				t.Error("PickByIndices with out-of-bounds index: expected panic, got none")
			}
		}()

		generic.PickByIndices(source, []int{99})
	})
}

func TestExpandHomeHermetic(t *testing.T) {
	// t.Setenv is incompatible with t.Parallel: the env mutation must not race
	// with other parallel tests.  Run this test (and its subtests) serially.

	// Use t.Setenv to override HOME so the test is independent of the real home
	// directory and produces a fully deterministic expected value.
	fakeHome := t.TempDir()
	t.Setenv("HOME", fakeHome)

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "tilde-slash prefix expands to fake home",
			input: "~/projects/foo",
			want:  filepath.Join(fakeHome, "projects/foo"),
		},
		{
			name:  "tilde alone expands to fake home",
			input: "~",
			want:  fakeHome,
		},
		{
			name:  "no tilde is returned unchanged",
			input: "/absolute/path",
			want:  "/absolute/path",
		},
		{
			name:  "empty string is returned unchanged",
			input: "",
			want:  "",
		},
	}

	for _, tc := range tests { //nolint:paralleltest // t.Setenv is incompatible with t.Parallel
		t.Run(tc.name, func(t *testing.T) {
			got := generic.ExpandHome(tc.input)
			if got != tc.want {
				t.Errorf("ExpandHome(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestIsURLMailto(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			// mailto: has a scheme but no host component — url.Parse gives an
			// opaque URI, so u.Host is empty. IsURL must return false.
			name:  "mailto scheme has no host",
			input: "mailto:user@example.com",
			want:  false,
		},
		{
			// A bare path segment with a colon looks like a scheme to url.Parse
			// only in specific circumstances; plain "host:path" without "//" is
			// treated as a relative URL with no scheme by url.Parse.
			name:  "host-colon-path is not a URL",
			input: "example.com:8080",
			want:  false,
		},
		{
			// urn: is a valid scheme but URNs have no host.
			name:  "urn scheme has no host",
			input: "urn:isbn:0451450523",
			want:  false,
		},
		{
			// data: URIs have a scheme but never a host.
			name:  "data URI has no host",
			input: "data:text/plain;base64,SGVsbG8=",
			want:  false,
		},
		{
			// Confirm a normal https URL still returns true (regression guard).
			name:  "https URL with path is valid",
			input: "https://example.com/path?q=1",
			want:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := generic.IsURL(tc.input)
			if got != tc.want {
				t.Errorf("IsURL(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}
