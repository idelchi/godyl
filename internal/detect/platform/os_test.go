package platform_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/idelchi/godyl/internal/detect/platform"
)

func TestOSParseFrom(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         string
		wantName      string
		wantCanonical string
		wantErrIs     error
	}{
		{
			name:          "linux exact",
			input:         "linux",
			wantName:      "linux",
			wantCanonical: "linux",
		},
		{
			name:          "darwin exact",
			input:         "darwin",
			wantName:      "darwin",
			wantCanonical: "darwin",
		},
		{
			name:          "macos alias for darwin",
			input:         "macos",
			wantName:      "macos",
			wantCanonical: "darwin",
		},
		{
			name:          "windows exact",
			input:         "windows",
			wantName:      "windows",
			wantCanonical: "windows",
		},
		{
			name:          "freebsd exact",
			input:         "freebsd",
			wantName:      "freebsd",
			wantCanonical: "freebsd",
		},
		{
			name:          "android exact",
			input:         "android",
			wantName:      "android",
			wantCanonical: "android",
		},
		{
			name:          "netbsd exact",
			input:         "netbsd",
			wantName:      "netbsd",
			wantCanonical: "netbsd",
		},
		{
			name:          "openbsd exact",
			input:         "openbsd",
			wantName:      "openbsd",
			wantCanonical: "openbsd",
		},
		{
			name:      "unknown OS returns ErrParse",
			input:     "unknown_os",
			wantErrIs: platform.ErrParse,
		},
		{
			name:      "empty string returns error",
			input:     "",
			wantErrIs: platform.ErrParse,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var o platform.OS

			err := o.ParseFrom(tc.input, strings.EqualFold, strings.Contains)

			if tc.wantErrIs != nil {
				if err == nil {
					t.Fatalf("ParseFrom(%q) = nil, want error", tc.input)
				}

				if !errors.Is(err, tc.wantErrIs) {
					t.Fatalf("ParseFrom(%q) error = %v, want errors.Is %v", tc.input, err, tc.wantErrIs)
				}

				return
			}

			if err != nil {
				t.Fatalf("ParseFrom(%q) unexpected error: %v", tc.input, err)
			}

			if o.Name != tc.wantName {
				t.Errorf("ParseFrom(%q): Name = %q, want %q", tc.input, o.Name, tc.wantName)
			}

			if o.Type() != tc.wantCanonical {
				t.Errorf("ParseFrom(%q): Type() = %q, want %q", tc.input, o.Type(), tc.wantCanonical)
			}
		})
	}
}

func TestOSParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         string
		wantCanonical string
		wantErr       bool
	}{
		{name: "linux parses via Parse()", input: "linux", wantCanonical: "linux"},
		{name: "darwin parses via Parse()", input: "darwin", wantCanonical: "darwin"},
		{name: "unknown fails via Parse()", input: "unknownos", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			o := platform.OS{Name: tc.input}

			err := o.Parse()

			if tc.wantErr {
				if err == nil {
					t.Fatalf("Parse() expected error, got nil")
				}

				return
			}

			if err != nil {
				t.Fatalf("Parse() unexpected error: %v", err)
			}

			if o.Type() != tc.wantCanonical {
				t.Errorf("Parse() Type() = %q, want %q", o.Type(), tc.wantCanonical)
			}
		})
	}
}

func TestOSIsCompatibleWith(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		left  string
		right string
		want  bool
	}{
		{
			name:  "linux compatible with linux",
			left:  "linux",
			right: "linux",
			want:  true,
		},
		{
			name:  "linux not compatible with darwin",
			left:  "linux",
			right: "darwin",
			want:  false,
		},
		{
			name:  "darwin compatible with darwin",
			left:  "darwin",
			right: "darwin",
			want:  true,
		},
		{
			name:  "macos compatible with darwin (same canonical)",
			left:  "macos",
			right: "darwin",
			want:  true,
		},
		{
			name:  "darwin compatible with macos (same canonical)",
			left:  "darwin",
			right: "macos",
			want:  true,
		},
		{
			name:  "windows not compatible with linux",
			left:  "windows",
			right: "linux",
			want:  false,
		},
		{
			name:  "freebsd not compatible with darwin",
			left:  "freebsd",
			right: "darwin",
			want:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var left, right platform.OS

			if err := left.ParseFrom(tc.left, strings.EqualFold); err != nil {
				t.Fatalf("ParseFrom(%q) unexpected error: %v", tc.left, err)
			}

			if err := right.ParseFrom(tc.right, strings.EqualFold); err != nil {
				t.Fatalf("ParseFrom(%q) unexpected error: %v", tc.right, err)
			}

			got := left.IsCompatibleWith(right)
			if got != tc.want {
				t.Errorf("IsCompatibleWith: %q vs %q = %v, want %v", tc.left, tc.right, got, tc.want)
			}
		})
	}
}

func TestOSIsCompatibleWithZeroValue(t *testing.T) {
	t.Parallel()

	var zero, linux platform.OS

	if err := linux.ParseFrom("linux", strings.EqualFold); err != nil {
		t.Fatalf("ParseFrom unexpected error: %v", err)
	}

	if zero.IsCompatibleWith(linux) {
		t.Error("zero-value OS should not be compatible with linux")
	}

	if linux.IsCompatibleWith(zero) {
		t.Error("linux should not be compatible with zero-value OS")
	}

	if zero.IsCompatibleWith(zero) {
		t.Error("zero-value OS should not be compatible with itself")
	}
}

func TestOSIs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		left  string
		right string
		want  bool
	}{
		{
			name:  "linux matches linux (same alias)",
			left:  "linux",
			right: "linux",
			want:  true,
		},
		{
			name:  "darwin matches darwin (same alias)",
			left:  "darwin",
			right: "darwin",
			want:  true,
		},
		{
			name:  "linux does not match darwin",
			left:  "linux",
			right: "darwin",
			want:  false,
		},
		{
			name:  "darwin does not match macos (different alias, same canonical)",
			left:  "darwin",
			right: "macos",
			want:  false,
		},
		{
			name:  "macos matches macos (same alias)",
			left:  "macos",
			right: "macos",
			want:  true,
		},
		{
			name:  "windows matches windows (same alias)",
			left:  "windows",
			right: "windows",
			want:  true,
		},
		{
			name:  "windows does not match win (different alias, same canonical)",
			left:  "windows",
			right: "win",
			want:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var left, right platform.OS

			if err := left.ParseFrom(tc.left, strings.EqualFold); err != nil {
				t.Fatalf("ParseFrom(%q) unexpected error: %v", tc.left, err)
			}

			if err := right.ParseFrom(tc.right, strings.EqualFold); err != nil {
				t.Fatalf("ParseFrom(%q) unexpected error: %v", tc.right, err)
			}

			got := left.Is(right)
			if got != tc.want {
				t.Errorf("Is: %q vs %q = %v, want %v", tc.left, tc.right, got, tc.want)
			}
		})
	}
}

func TestOSIsNil(t *testing.T) {
	t.Parallel()

	t.Run("zero-value OS is nil", func(t *testing.T) {
		t.Parallel()

		var o platform.OS

		if !o.IsNil() {
			t.Error("zero-value OS.IsNil() = false, want true")
		}
	})

	t.Run("parsed OS is not nil", func(t *testing.T) {
		t.Parallel()

		var o platform.OS

		if err := o.ParseFrom("linux", strings.EqualFold); err != nil {
			t.Fatalf("ParseFrom(%q) unexpected error: %v", "linux", err)
		}

		if o.IsNil() {
			t.Error("parsed OS.IsNil() = true, want false")
		}
	})

	t.Run("darwin parsed OS is not nil", func(t *testing.T) {
		t.Parallel()

		var o platform.OS

		if err := o.ParseFrom("darwin", strings.EqualFold); err != nil {
			t.Fatalf("ParseFrom(%q) unexpected error: %v", "darwin", err)
		}

		if o.IsNil() {
			t.Error("darwin OS.IsNil() = true, want false")
		}
	})
}

func TestOSParseFromAliases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         string
		wantCanonical string
	}{
		{
			name:          "mac alias maps to darwin",
			input:         "mac",
			wantCanonical: "darwin",
		},
		{
			name:          "osx alias maps to darwin",
			input:         "osx",
			wantCanonical: "darwin",
		},
		{
			name:          "win alias maps to windows",
			input:         "win",
			wantCanonical: "windows",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var o platform.OS

			if err := o.ParseFrom(tc.input, strings.EqualFold); err != nil {
				t.Fatalf("ParseFrom(%q) unexpected error: %v", tc.input, err)
			}

			if o.Type() != tc.wantCanonical {
				t.Errorf("ParseFrom(%q): Type() = %q, want %q", tc.input, o.Type(), tc.wantCanonical)
			}
		})
	}
}

func TestOSStringNetBSD(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "netbsd string",
			input: "netbsd",
			want:  "netbsd",
		},
		{
			name:  "openbsd string",
			input: "openbsd",
			want:  "openbsd",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var o platform.OS

			if err := o.ParseFrom(tc.input, strings.EqualFold); err != nil {
				t.Fatalf("ParseFrom(%q) unexpected error: %v", tc.input, err)
			}

			got := o.String()
			if got != tc.want {
				t.Errorf("String() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestOSString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "linux string",
			input: "linux",
			want:  "linux",
		},
		{
			name:  "darwin string",
			input: "darwin",
			want:  "darwin",
		},
		{
			name:  "macos alias string is canonical darwin",
			input: "macos",
			want:  "darwin",
		},
		{
			name:  "windows string",
			input: "windows",
			want:  "windows",
		},
		{
			name:  "freebsd string",
			input: "freebsd",
			want:  "freebsd",
		},
		{
			name:  "android string",
			input: "android",
			want:  "android",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var o platform.OS

			if err := o.ParseFrom(tc.input, strings.EqualFold); err != nil {
				t.Fatalf("ParseFrom(%q) unexpected error: %v", tc.input, err)
			}

			got := o.String()
			if got != tc.want {
				t.Errorf("String() = %q, want %q", got, tc.want)
			}
		})
	}
}
