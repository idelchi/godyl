package platform_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/idelchi/godyl/internal/detect/platform"
)

func TestLibraryParseFrom(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         string
		wantName      string
		wantCanonical string
		wantErrIs     error
	}{
		{
			name:          "musl exact",
			input:         "musl",
			wantName:      "musl",
			wantCanonical: "musl",
		},
		{
			name:          "gnu exact",
			input:         "gnu",
			wantName:      "gnu",
			wantCanonical: "gnu",
		},
		{
			name:          "glibc alias for gnu",
			input:         "glibc",
			wantName:      "glibc",
			wantCanonical: "gnu",
		},
		{
			name:          "msvc exact",
			input:         "msvc",
			wantName:      "msvc",
			wantCanonical: "msvc",
		},
		{
			name:          "android exact",
			input:         "android",
			wantName:      "android",
			wantCanonical: "android",
		},
		{
			name:      "unknown library returns ErrParse",
			input:     "unknown_lib",
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

			var l platform.Library

			err := l.ParseFrom(tc.input, strings.EqualFold, strings.Contains)

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

			if l.Name != tc.wantName {
				t.Errorf("ParseFrom(%q): Name = %q, want %q", tc.input, l.Name, tc.wantName)
			}

			if l.String() != tc.wantCanonical {
				t.Errorf("ParseFrom(%q): String() = %q, want %q", tc.input, l.String(), tc.wantCanonical)
			}
		})
	}
}

func TestLibraryParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         string
		wantCanonical string
		wantErr       bool
	}{
		{name: "gnu parses via Parse()", input: "gnu", wantCanonical: "gnu"},
		{name: "musl parses via Parse()", input: "musl", wantCanonical: "musl"},
		{name: "unknown fails via Parse()", input: "unknownlib", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			l := platform.Library{Name: tc.input}

			err := l.Parse()

			if tc.wantErr {
				if err == nil {
					t.Fatalf("Parse() expected error, got nil")
				}

				return
			}

			if err != nil {
				t.Fatalf("Parse() unexpected error: %v", err)
			}

			if l.String() != tc.wantCanonical {
				t.Errorf("Parse() String() = %q, want %q", l.String(), tc.wantCanonical)
			}
		})
	}
}

func TestLibraryIsCompatibleWith(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		left  string
		right string
		want  bool
	}{
		{
			name:  "gnu compatible with gnu",
			left:  "gnu",
			right: "gnu",
			want:  true,
		},
		{
			name:  "musl compatible with musl",
			left:  "musl",
			right: "musl",
			want:  true,
		},
		{
			name:  "gnu compatible with musl per matrix",
			left:  "gnu",
			right: "musl",
			want:  true,
		},
		{
			name:  "gnu not compatible with android",
			left:  "gnu",
			right: "android",
			want:  false,
		},
		{
			name:  "android not compatible with gnu",
			left:  "android",
			right: "gnu",
			want:  false,
		},
		{
			name:  "msvc compatible with msvc",
			left:  "msvc",
			right: "msvc",
			want:  true,
		},
		{
			name:  "msvc compatible with gnu per matrix",
			left:  "msvc",
			right: "gnu",
			want:  true,
		},
		{
			name:  "msvc not compatible with musl",
			left:  "msvc",
			right: "musl",
			want:  false,
		},
		{
			name:  "libSystem compatible with libSystem",
			left:  "libSystem",
			right: "libSystem",
			want:  true,
		},
		{
			name:  "libSystem compatible with gnu per matrix",
			left:  "libSystem",
			right: "gnu",
			want:  true,
		},
		{
			name:  "libSystem compatible with musl per matrix",
			left:  "libSystem",
			right: "musl",
			want:  true,
		},
		{
			name:  "libSystem not compatible with android",
			left:  "libSystem",
			right: "android",
			want:  false,
		},
		{
			name:  "musl not compatible with android",
			left:  "musl",
			right: "android",
			want:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var left, right platform.Library

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

func TestLibraryIsCompatibleWithZeroValue(t *testing.T) {
	t.Parallel()

	var zero, gnu platform.Library

	if err := gnu.ParseFrom("gnu", strings.EqualFold); err != nil {
		t.Fatalf("ParseFrom unexpected error: %v", err)
	}

	if zero.IsCompatibleWith(gnu) {
		t.Error("zero-value Library should not be compatible with gnu")
	}

	if gnu.IsCompatibleWith(zero) {
		t.Error("gnu should not be compatible with zero-value Library")
	}

	if zero.IsCompatibleWith(zero) {
		t.Error("zero-value Library should not be compatible with itself")
	}
}

func TestLibraryIs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		left  string
		right string
		want  bool
	}{
		{
			name:  "gnu matches gnu (same alias)",
			left:  "gnu",
			right: "gnu",
			want:  true,
		},
		{
			name:  "musl matches musl (same alias)",
			left:  "musl",
			right: "musl",
			want:  true,
		},
		{
			name:  "gnu does not match musl",
			left:  "gnu",
			right: "musl",
			want:  false,
		},
		{
			name:  "gnu does not match glibc (different alias, same canonical)",
			left:  "gnu",
			right: "glibc",
			want:  false,
		},
		{
			name:  "glibc matches glibc (same alias)",
			left:  "glibc",
			right: "glibc",
			want:  true,
		},
		{
			name:  "msvc matches msvc (same alias)",
			left:  "msvc",
			right: "msvc",
			want:  true,
		},
		{
			name:  "msvc does not match gnu",
			left:  "msvc",
			right: "gnu",
			want:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var left, right platform.Library

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

func TestLibraryIsNil(t *testing.T) {
	t.Parallel()

	t.Run("zero-value library is nil", func(t *testing.T) {
		t.Parallel()

		var l platform.Library

		if !l.IsNil() {
			t.Error("zero-value Library.IsNil() = false, want true")
		}
	})

	t.Run("gnu parsed library is not nil", func(t *testing.T) {
		t.Parallel()

		var l platform.Library

		if err := l.ParseFrom("gnu", strings.EqualFold); err != nil {
			t.Fatalf("ParseFrom(%q) unexpected error: %v", "gnu", err)
		}

		if l.IsNil() {
			t.Error("gnu Library.IsNil() = true, want false")
		}
	})

	t.Run("musl parsed library is not nil", func(t *testing.T) {
		t.Parallel()

		var l platform.Library

		if err := l.ParseFrom("musl", strings.EqualFold); err != nil {
			t.Fatalf("ParseFrom(%q) unexpected error: %v", "musl", err)
		}

		if l.IsNil() {
			t.Error("musl Library.IsNil() = true, want false")
		}
	})
}

func TestLibraryDefault(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		osName     string
		distroName string
		wantName   string
	}{
		{
			name:       "linux with debian distro returns gnu",
			osName:     "linux",
			distroName: "debian",
			wantName:   "gnu",
		},
		{
			name:       "linux with ubuntu distro returns gnu",
			osName:     "linux",
			distroName: "ubuntu",
			wantName:   "gnu",
		},
		{
			name:       "linux with alpine distro returns musl",
			osName:     "linux",
			distroName: "alpine",
			wantName:   "musl",
		},
		{
			name:       "linux with unknown distro returns gnu",
			osName:     "linux",
			distroName: "",
			wantName:   "gnu",
		},
		{
			name:     "darwin returns libSystem",
			osName:   "darwin",
			wantName: "libSystem",
		},
		{
			name:     "windows returns msvc",
			osName:   "windows",
			wantName: "msvc",
		},
		{
			name:     "android returns android",
			osName:   "android",
			wantName: "android",
		},
		{
			name:     "freebsd returns empty Library (fallthrough)",
			osName:   "freebsd",
			wantName: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var o platform.OS

			if err := o.ParseFrom(tc.osName, strings.EqualFold); err != nil {
				t.Fatalf("OS.ParseFrom(%q) unexpected error: %v", tc.osName, err)
			}

			var distro platform.Distribution

			if tc.distroName != "" {
				if err := distro.ParseFrom(tc.distroName, strings.EqualFold); err != nil {
					t.Fatalf("Distribution.ParseFrom(%q) unexpected error: %v", tc.distroName, err)
				}
			}

			var l platform.Library

			got := l.Default(o, distro)

			if got.Name != tc.wantName {
				t.Errorf("Default(%q, %q): Name = %q, want %q", tc.osName, tc.distroName, got.Name, tc.wantName)
			}
		})
	}
}
