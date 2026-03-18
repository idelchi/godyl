package platform_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/idelchi/godyl/internal/detect/platform"
)

func TestArchParseFrom(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         string
		wantName      string
		wantCanonical string
		wantErrIs     error
	}{
		{
			name:          "amd64 exact",
			input:         "amd64",
			wantName:      "amd64",
			wantCanonical: "amd64",
		},
		{
			name:          "x86_64 alias for amd64",
			input:         "x86_64",
			wantName:      "x86_64",
			wantCanonical: "amd64",
		},
		{
			name:          "aarch64 alias for arm64",
			input:         "aarch64",
			wantName:      "aarch64",
			wantCanonical: "arm64",
		},
		{
			name:          "arm64 exact",
			input:         "arm64",
			wantName:      "arm64",
			wantCanonical: "arm64",
		},
		{
			name:          "i386 alias for 386",
			input:         "i386",
			wantName:      "i386",
			wantCanonical: "386",
		},
		{
			name:          "i686 alias for 386",
			input:         "i686",
			wantName:      "i686",
			wantCanonical: "386",
		},
		{
			name:      "unknown architecture returns ErrParse",
			input:     "unknown_arch",
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

			var a platform.Architecture

			err := a.ParseFrom(tc.input, strings.EqualFold, strings.Contains)

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

			if a.Name != tc.wantName {
				t.Errorf("ParseFrom(%q): Name = %q, want %q", tc.input, a.Name, tc.wantName)
			}

			if a.Type() != tc.wantCanonical {
				t.Errorf("ParseFrom(%q): Type() = %q, want %q", tc.input, a.Type(), tc.wantCanonical)
			}
		})
	}
}

func TestArchParseFromContains(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         string
		wantCanonical string
	}{
		{
			name:          "linux-amd64 matches amd64 via Contains",
			input:         "linux-amd64",
			wantCanonical: "amd64",
		},
		{
			name:          "x86_64-linux-gnu matches amd64 via Contains",
			input:         "x86_64-linux-gnu",
			wantCanonical: "amd64",
		},
		{
			name:          "aarch64-linux-gnu matches arm64 via Contains",
			input:         "aarch64-linux-gnu",
			wantCanonical: "arm64",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var a platform.Architecture

			if err := a.ParseFrom(tc.input, strings.Contains); err != nil {
				t.Fatalf("ParseFrom(%q) unexpected error: %v", tc.input, err)
			}

			if a.Type() != tc.wantCanonical {
				t.Errorf("ParseFrom(%q): Type() = %q, want %q", tc.input, a.Type(), tc.wantCanonical)
			}
		})
	}
}

func TestArchIsCompatibleWith(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		left  string
		right string
		want  bool
	}{
		{
			name:  "amd64 compatible with amd64",
			left:  "amd64",
			right: "amd64",
			want:  true,
		},
		{
			name:  "arm64 compatible with arm64",
			left:  "arm64",
			right: "arm64",
			want:  true,
		},
		{
			name:  "amd64 not compatible with arm64",
			left:  "amd64",
			right: "arm64",
			want:  false,
		},
		{
			name:  "armv7 compatible with armv5 (higher version runs lower)",
			left:  "armv7",
			right: "armv5",
			want:  true,
		},
		{
			name:  "armv5 not compatible with armv7 (lower version cannot run higher)",
			left:  "armv5",
			right: "armv7",
			want:  false,
		},
		{
			name:  "arm not compatible with arm64 (different canonical)",
			left:  "armhf",
			right: "arm64",
			want:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var left, right platform.Architecture

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

func TestArchIsCompatibleWithZeroValue(t *testing.T) {
	t.Parallel()

	var zero, amd64 platform.Architecture

	if err := amd64.ParseFrom("amd64", strings.EqualFold); err != nil {
		t.Fatalf("ParseFrom unexpected error: %v", err)
	}

	if zero.IsCompatibleWith(amd64) {
		t.Error("zero-value Architecture should not be compatible with amd64")
	}

	if amd64.IsCompatibleWith(zero) {
		t.Error("amd64 should not be compatible with zero-value Architecture")
	}

	if zero.IsCompatibleWith(zero) {
		t.Error("zero-value Architecture should not be compatible with itself")
	}
}

func TestArchString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "amd64 string",
			input: "amd64",
			want:  "amd64",
		},
		{
			name:  "x86_64 alias string is canonical amd64",
			input: "x86_64",
			want:  "amd64",
		},
		{
			name:  "arm64 string",
			input: "arm64",
			want:  "arm64",
		},
		{
			name:  "aarch64 alias string is canonical arm64",
			input: "aarch64",
			want:  "arm64",
		},
		{
			name:  "armv7 string includes version",
			input: "armv7",
			want:  "armv7",
		},
		{
			name:  "armv5 string includes version",
			input: "armv5",
			want:  "armv5",
		},
		{
			name:  "armhf string is armv7",
			input: "armhf",
			want:  "armv7",
		},
		{
			name:  "386 string",
			input: "386",
			want:  "386",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var a platform.Architecture

			if err := a.ParseFrom(tc.input, strings.EqualFold); err != nil {
				t.Fatalf("ParseFrom(%q) unexpected error: %v", tc.input, err)
			}

			got := a.String()
			if got != tc.want {
				t.Errorf("String() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestArchIs64Bit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "amd64 is 64-bit", input: "amd64", want: true},
		{name: "arm64 is 64-bit", input: "arm64", want: true},
		{name: "386 is not 64-bit", input: "386", want: false},
		{name: "armv7 is not 64-bit", input: "armv7", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var a platform.Architecture

			if err := a.ParseFrom(tc.input, strings.EqualFold); err != nil {
				t.Fatalf("ParseFrom(%q) unexpected error: %v", tc.input, err)
			}

			got := a.Is64Bit()
			if got != tc.want {
				t.Errorf("Is64Bit() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestArchIsARM(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "arm64 is ARM", input: "arm64", want: true},
		{name: "armv7 is ARM", input: "armv7", want: true},
		{name: "armv5 is ARM", input: "armv5", want: true},
		{name: "amd64 is not ARM", input: "amd64", want: false},
		{name: "386 is not ARM", input: "386", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var a platform.Architecture

			if err := a.ParseFrom(tc.input, strings.EqualFold); err != nil {
				t.Fatalf("ParseFrom(%q) unexpected error: %v", tc.input, err)
			}

			got := a.IsARM()
			if got != tc.want {
				t.Errorf("IsARM() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestArchIsX86(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "amd64 is x86", input: "amd64", want: true},
		{name: "386 is x86", input: "386", want: true},
		{name: "arm64 is not x86", input: "arm64", want: false},
		{name: "armv7 is not x86", input: "armv7", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var a platform.Architecture

			if err := a.ParseFrom(tc.input, strings.EqualFold); err != nil {
				t.Fatalf("ParseFrom(%q) unexpected error: %v", tc.input, err)
			}

			got := a.IsX86()
			if got != tc.want {
				t.Errorf("IsX86() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestArchParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         string
		wantCanonical string
		wantErr       bool
	}{
		{name: "amd64 parses via Parse()", input: "amd64", wantCanonical: "amd64"},
		{name: "arm64 parses via Parse()", input: "arm64", wantCanonical: "arm64"},
		{name: "unknown fails via Parse()", input: "unknownarch", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			a := platform.Architecture{Name: tc.input}

			err := a.Parse()

			if tc.wantErr {
				if err == nil {
					t.Fatalf("Parse() expected error, got nil")
				}

				return
			}

			if err != nil {
				t.Fatalf("Parse() unexpected error: %v", err)
			}

			if a.Type() != tc.wantCanonical {
				t.Errorf("Parse() Type() = %q, want %q", a.Type(), tc.wantCanonical)
			}
		})
	}
}

func TestArchIs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		left  string
		right string
		want  bool
	}{
		{
			name:  "amd64 matches amd64 (same alias)",
			left:  "amd64",
			right: "amd64",
			want:  true,
		},
		{
			name:  "arm64 matches arm64 (same alias)",
			left:  "arm64",
			right: "arm64",
			want:  true,
		},
		{
			name:  "armv7 matches armv7 (same alias)",
			left:  "armv7",
			right: "armv7",
			want:  true,
		},
		{
			name:  "amd64 does not match arm64",
			left:  "amd64",
			right: "arm64",
			want:  false,
		},
		{
			name:  "amd64 does not match x86_64 (different alias, same canonical)",
			left:  "amd64",
			right: "x86_64",
			want:  false,
		},
		{
			name:  "armv5 does not match armv7 (different alias)",
			left:  "armv5",
			right: "armv7",
			want:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var left, right platform.Architecture

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

func TestArchAliases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         string
		wantAliases   []string
		wantNilResult bool
	}{
		{
			name:        "amd64 aliases",
			input:       "amd64",
			wantAliases: []string{"x86_64", "x64", "win64"},
		},
		{
			name:        "arm64 aliases",
			input:       "arm64",
			wantAliases: []string{"aarch64"},
		},
		{
			name:        "386 aliases",
			input:       "386",
			wantAliases: []string{"amd32", "x86", "i386", "i686", "win32"},
		},
		{
			name:        "arm aliases include armv variants",
			input:       "armv7",
			wantAliases: []string{"armv7", "armv6", "armv5", "armel", "armhf", "arm"},
		},
		{
			name:          "zero-value returns nil aliases",
			input:         "",
			wantNilResult: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var a platform.Architecture

			if tc.input != "" {
				if err := a.ParseFrom(tc.input, strings.EqualFold); err != nil {
					t.Fatalf("ParseFrom(%q) unexpected error: %v", tc.input, err)
				}
			}

			got := a.Aliases()

			if tc.wantNilResult {
				if got != nil {
					t.Errorf("Aliases() = %v, want nil", got)
				}

				return
			}

			if len(got) != len(tc.wantAliases) {
				t.Fatalf(
					"Aliases() = %v (len %d), want %v (len %d)",
					got,
					len(got),
					tc.wantAliases,
					len(tc.wantAliases),
				)
			}

			for i, alias := range tc.wantAliases {
				if got[i] != alias {
					t.Errorf("Aliases()[%d] = %q, want %q", i, got[i], alias)
				}
			}
		})
	}
}

func TestArchIsNil(t *testing.T) {
	t.Parallel()

	t.Run("zero-value architecture is nil", func(t *testing.T) {
		t.Parallel()

		var a platform.Architecture

		if !a.IsNil() {
			t.Error("zero-value Architecture.IsNil() = false, want true")
		}
	})

	t.Run("parsed architecture is not nil", func(t *testing.T) {
		t.Parallel()

		var a platform.Architecture

		if err := a.ParseFrom("amd64", strings.EqualFold); err != nil {
			t.Fatalf("ParseFrom(%q) unexpected error: %v", "amd64", err)
		}

		if a.IsNil() {
			t.Error("parsed Architecture.IsNil() = true, want false")
		}
	})

	t.Run("arm64 parsed architecture is not nil", func(t *testing.T) {
		t.Parallel()

		var a platform.Architecture

		if err := a.ParseFrom("arm64", strings.EqualFold); err != nil {
			t.Fatalf("ParseFrom(%q) unexpected error: %v", "arm64", err)
		}

		if a.IsNil() {
			t.Error("arm64 Architecture.IsNil() = true, want false")
		}
	})
}

func TestArchParseFromMixedCase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         string
		wantCanonical string
	}{
		{
			name:          "AMD64 uppercase matches amd64",
			input:         "AMD64",
			wantCanonical: "amd64",
		},
		{
			name:          "X86_64 uppercase matches amd64",
			input:         "X86_64",
			wantCanonical: "amd64",
		},
		{
			name:          "ARM64 uppercase matches arm64",
			input:         "ARM64",
			wantCanonical: "arm64",
		},
		{
			name:          "AARCH64 uppercase matches arm64",
			input:         "AARCH64",
			wantCanonical: "arm64",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var a platform.Architecture

			if err := a.ParseFrom(tc.input, strings.EqualFold); err != nil {
				t.Fatalf("ParseFrom(%q) unexpected error: %v", tc.input, err)
			}

			if a.Type() != tc.wantCanonical {
				t.Errorf("ParseFrom(%q): Type() = %q, want %q", tc.input, a.Type(), tc.wantCanonical)
			}
		})
	}
}

func TestArchTo32BitUserLand(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         string
		wantCanonical string
		wantString    string
	}{
		{
			name:          "amd64 converts to 386",
			input:         "amd64",
			wantCanonical: "386",
			wantString:    "386",
		},
		{
			name:          "arm64 converts to armv7",
			input:         "arm64",
			wantCanonical: "arm",
			wantString:    "armv7",
		},
		{
			name:          "386 unchanged",
			input:         "386",
			wantCanonical: "386",
			wantString:    "386",
		},
		{
			name:          "armv7 unchanged",
			input:         "armv7",
			wantCanonical: "arm",
			wantString:    "armv7",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var a platform.Architecture

			if err := a.ParseFrom(tc.input, strings.EqualFold); err != nil {
				t.Fatalf("ParseFrom(%q) unexpected error: %v", tc.input, err)
			}

			a.To32BitUserLand()

			if a.Type() != tc.wantCanonical {
				t.Errorf("To32BitUserLand() Type() = %q, want %q", a.Type(), tc.wantCanonical)
			}

			if a.String() != tc.wantString {
				t.Errorf("To32BitUserLand() String() = %q, want %q", a.String(), tc.wantString)
			}
		})
	}
}
