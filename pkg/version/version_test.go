package version_test

import (
	"testing"

	"github.com/idelchi/godyl/pkg/version"
)

func TestParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string // empty string means nil is expected
	}{
		{name: "prefixed v", input: "v1.2.3", want: "1.2.3"},
		{name: "no prefix", input: "1.2.3", want: "1.2.3"},
		{name: "tool prefix with v and arch suffix", input: "tool-v1.2.3-linux-amd64", want: "1.2.3-linux-amd64"},
		{name: "tool prefix no v with os suffix", input: "tool-1.2.3-linux", want: "1.2.3-linux"},
		{name: "prerelease with v", input: "v1.2.3-beta.1", want: "1.2.3-beta.1"},
		{name: "empty string", input: "", want: ""},
		{name: "no version present", input: "notaversion", want: ""},
		{name: "tool prefix with v only", input: "tool-v2.0.0", want: "2.0.0"},
		{name: "prerelease suffix after v", input: "v0.1.0-rest", want: "0.1.0-rest"},
		{
			name:  "version embedded mid-string with extra suffix",
			input: "prefix-1.2.3-suffix-extra",
			want:  "1.2.3-suffix-extra",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := version.Parse(tc.input)

			if tc.want == "" {
				if got != nil {
					t.Errorf("Parse(%q) = %q, want nil", tc.input, got.String())
				}

				return
			}

			if got == nil {
				t.Fatalf("Parse(%q) = nil, want %q", tc.input, tc.want)
			}

			if got.String() != tc.want {
				t.Errorf("Parse(%q).String() = %q, want %q", tc.input, got.String(), tc.want)
			}
		})
	}
}

func TestEqual(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		a, b string
		want bool
	}{
		{name: "identical with v prefix", a: "v1.2.3", b: "v1.2.3", want: true},
		{name: "v prefix vs bare", a: "v1.2.3", b: "1.2.3", want: true},
		{name: "different patch", a: "v1.2.3", b: "v1.2.4", want: false},
		{name: "invalid a", a: "bad", b: "v1.0.0", want: false},
		{name: "invalid b", a: "v1.0.0", b: "bad", want: false},
		{name: "both invalid returns false (nil comparison)", a: "bad", b: "bad", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := version.Equal(tc.a, tc.b)
			if got != tc.want {
				t.Errorf("Equal(%q, %q) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestParseZeroVersion(t *testing.T) {
	t.Parallel()

	// "0.0.0" is a valid semantic version; Parse must return a non-nil result.
	got := version.Parse("0.0.0")
	if got == nil {
		t.Fatal("Parse(\"0.0.0\") = nil, want non-nil *semver.Version")
	}

	const want = "0.0.0"

	if got.String() != want {
		t.Errorf("Parse(\"0.0.0\").String() = %q, want %q", got.String(), want)
	}
}

func TestLessThan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		a, b string
		want bool
	}{
		{name: "less than", a: "v1.2.3", b: "v1.2.4", want: true},
		{name: "greater than", a: "v1.2.4", b: "v1.2.3", want: false},
		{name: "equal", a: "v1.2.3", b: "v1.2.3", want: false},
		{name: "invalid a or b causes AnyNil fallback returns true", a: "bad", b: "v1.0.0", want: true},
		{name: "both invalid returns true (failure mode)", a: "bad", b: "bad", want: true},
		{name: "valid a invalid b returns true", a: "v1.0.0", b: "bad", want: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := version.LessThan(tc.a, tc.b)
			if got != tc.want {
				t.Errorf("LessThan(%q, %q) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestLessThanPreReleaseVsRelease(t *testing.T) {
	t.Parallel()

	// Per semver spec, a pre-release version has lower precedence than the
	// associated normal version: v1.0.0-alpha < v1.0.0.
	if !version.LessThan("v1.0.0-alpha", "v1.0.0") {
		t.Error("LessThan(\"v1.0.0-alpha\", \"v1.0.0\") = false, want true (pre-release < release per semver)")
	}
}

func TestEqualPreRelease(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		a, b string
		want bool
	}{
		{
			// Two different pre-release labels on the same base version are not equal.
			name: "alpha vs beta are not equal",
			a:    "v1.0.0-alpha",
			b:    "v1.0.0-beta",
			want: false,
		},
		{
			// The same pre-release label is equal to itself.
			name: "same pre-release label is equal",
			a:    "v1.0.0-alpha",
			b:    "v1.0.0-alpha",
			want: true,
		},
		{
			// A pre-release version is not equal to its release counterpart.
			name: "pre-release is not equal to release",
			a:    "v1.0.0-alpha",
			b:    "v1.0.0",
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := version.Equal(tc.a, tc.b)
			if got != tc.want {
				t.Errorf("Equal(%q, %q) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}
