package platform_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/idelchi/godyl/internal/detect/platform"
)

func TestDistributionParseFrom(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         string
		wantName      string
		wantCanonical string
		wantErrIs     error
	}{
		{
			name:          "alpine exact",
			input:         "alpine",
			wantName:      "alpine",
			wantCanonical: "alpine",
		},
		{
			name:          "ubuntu exact",
			input:         "ubuntu",
			wantName:      "ubuntu",
			wantCanonical: "ubuntu",
		},
		{
			name:          "debian exact",
			input:         "debian",
			wantName:      "debian",
			wantCanonical: "debian",
		},
		{
			name:          "centos exact",
			input:         "centos",
			wantName:      "centos",
			wantCanonical: "centos",
		},
		{
			name:          "arch exact",
			input:         "arch",
			wantName:      "arch",
			wantCanonical: "arch",
		},
		{
			name:          "redhat exact",
			input:         "redhat",
			wantName:      "redhat",
			wantCanonical: "redhat",
		},
		{
			name:          "rhel alias for redhat",
			input:         "rhel",
			wantName:      "rhel",
			wantCanonical: "redhat",
		},
		{
			name:          "raspbian exact",
			input:         "raspbian",
			wantName:      "raspbian",
			wantCanonical: "raspbian",
		},
		{
			name:          "raspberry alias for raspbian",
			input:         "raspberry",
			wantName:      "raspberry",
			wantCanonical: "raspbian",
		},
		{
			name:      "unknown distro returns ErrParse",
			input:     "unknown_distro",
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

			var d platform.Distribution

			err := d.ParseFrom(tc.input, strings.EqualFold, strings.Contains)

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

			if d.Name != tc.wantName {
				t.Errorf("ParseFrom(%q): Name = %q, want %q", tc.input, d.Name, tc.wantName)
			}

			if d.String() != tc.wantCanonical {
				t.Errorf("ParseFrom(%q): String() = %q, want %q", tc.input, d.String(), tc.wantCanonical)
			}
		})
	}
}

func TestDistributionParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         string
		wantCanonical string
		wantErr       bool
	}{
		{name: "alpine parses via Parse()", input: "alpine", wantCanonical: "alpine"},
		{name: "ubuntu parses via Parse()", input: "ubuntu", wantCanonical: "ubuntu"},
		{name: "unknown fails via Parse()", input: "unknowndistro", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			d := platform.Distribution{Name: tc.input}

			err := d.Parse()

			if tc.wantErr {
				if err == nil {
					t.Fatalf("Parse() expected error, got nil")
				}

				return
			}

			if err != nil {
				t.Fatalf("Parse() unexpected error: %v", err)
			}

			if d.String() != tc.wantCanonical {
				t.Errorf("Parse() String() = %q, want %q", d.String(), tc.wantCanonical)
			}
		})
	}
}

func TestDistributionIsUnset(t *testing.T) {
	t.Parallel()

	t.Run("zero-value distribution is unset", func(t *testing.T) {
		t.Parallel()

		var d platform.Distribution

		if !d.IsUnset() {
			t.Error("zero-value Distribution.IsUnset() = false, want true")
		}
	})

	t.Run("parsed distribution is not unset", func(t *testing.T) {
		t.Parallel()

		var d platform.Distribution

		if err := d.ParseFrom("alpine", strings.EqualFold); err != nil {
			t.Fatalf("ParseFrom(%q) unexpected error: %v", "alpine", err)
		}

		if d.IsUnset() {
			t.Error("parsed Distribution.IsUnset() = true, want false")
		}
	})
}

func TestDistributionIs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		left  string
		right string
		want  bool
	}{
		{
			name:  "alpine matches alpine (same alias)",
			left:  "alpine",
			right: "alpine",
			want:  true,
		},
		{
			name:  "ubuntu does not match debian (different alias)",
			left:  "ubuntu",
			right: "debian",
			want:  false,
		},
		{
			name:  "redhat does not match rhel (different alias, same canonical)",
			left:  "redhat",
			right: "rhel",
			want:  false,
		},
		{
			name:  "rhel matches rhel (same alias)",
			left:  "rhel",
			right: "rhel",
			want:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var left, right platform.Distribution

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

func TestDistributionIsZeroValue(t *testing.T) {
	t.Parallel()

	var zero, alpine platform.Distribution

	if err := alpine.ParseFrom("alpine", strings.EqualFold); err != nil {
		t.Fatalf("ParseFrom unexpected error: %v", err)
	}

	if zero.Is(alpine) {
		t.Error("zero-value Distribution should not match alpine")
	}

	if alpine.Is(zero) {
		t.Error("alpine should not match zero-value Distribution")
	}

	if zero.Is(zero) {
		t.Error("zero-value Distribution should not match itself")
	}
}

func TestDistributionIsCompatibleWith(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		left  string
		right string
		want  bool
	}{
		{
			name:  "alpine compatible with alpine (same canonical)",
			left:  "alpine",
			right: "alpine",
			want:  true,
		},
		{
			name:  "ubuntu compatible with ubuntu (same canonical)",
			left:  "ubuntu",
			right: "ubuntu",
			want:  true,
		},
		{
			name:  "redhat compatible with rhel (same canonical)",
			left:  "redhat",
			right: "rhel",
			want:  true,
		},
		{
			name:  "rhel compatible with redhat (same canonical)",
			left:  "rhel",
			right: "redhat",
			want:  true,
		},
		{
			name:  "raspbian compatible with raspberry (same canonical)",
			left:  "raspbian",
			right: "raspberry",
			want:  true,
		},
		{
			name:  "alpine not compatible with ubuntu (different canonical)",
			left:  "alpine",
			right: "ubuntu",
			want:  false,
		},
		{
			name:  "debian not compatible with centos (different canonical)",
			left:  "debian",
			right: "centos",
			want:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var left, right platform.Distribution

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

func TestDistributionIsCompatibleWithZeroValue(t *testing.T) {
	t.Parallel()

	var zero, alpine platform.Distribution

	if err := alpine.ParseFrom("alpine", strings.EqualFold); err != nil {
		t.Fatalf("ParseFrom unexpected error: %v", err)
	}

	if zero.IsCompatibleWith(alpine) {
		t.Error("zero-value Distribution should not be compatible with alpine")
	}

	if alpine.IsCompatibleWith(zero) {
		t.Error("alpine should not be compatible with zero-value Distribution")
	}

	if zero.IsCompatibleWith(zero) {
		t.Error("zero-value Distribution should not be compatible with itself")
	}
}

func TestDistributionIsNil(t *testing.T) {
	t.Parallel()

	t.Run("zero-value distribution is nil", func(t *testing.T) {
		t.Parallel()

		var d platform.Distribution

		if !d.IsNil() {
			t.Error("zero-value Distribution.IsNil() = false, want true")
		}
	})

	t.Run("parsed distribution is not nil", func(t *testing.T) {
		t.Parallel()

		var d platform.Distribution

		if err := d.ParseFrom("ubuntu", strings.EqualFold); err != nil {
			t.Fatalf("ParseFrom(%q) unexpected error: %v", "ubuntu", err)
		}

		if d.IsNil() {
			t.Error("parsed Distribution.IsNil() = true, want false")
		}
	})

	t.Run("alpine parsed distribution is not nil", func(t *testing.T) {
		t.Parallel()

		var d platform.Distribution

		if err := d.ParseFrom("alpine", strings.EqualFold); err != nil {
			t.Fatalf("ParseFrom(%q) unexpected error: %v", "alpine", err)
		}

		if d.IsNil() {
			t.Error("alpine Distribution.IsNil() = true, want false")
		}
	})
}

func TestDistributionParseFromContains(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         string
		wantCanonical string
	}{
		{
			name:          "ubuntu-22.04 matches ubuntu via Contains",
			input:         "ubuntu-22.04",
			wantCanonical: "ubuntu",
		},
		{
			name:          "alpine-linux matches alpine via Contains",
			input:         "alpine-linux",
			wantCanonical: "alpine",
		},
		{
			name:          "debian-bullseye matches debian via Contains",
			input:         "debian-bullseye",
			wantCanonical: "debian",
		},
		{
			name:          "raspberry-pi matches raspbian via Contains on alias",
			input:         "raspberry-pi",
			wantCanonical: "raspbian",
		},
		{
			name:          "centos-stream matches centos via Contains",
			input:         "centos-stream",
			wantCanonical: "centos",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var d platform.Distribution

			if err := d.ParseFrom(tc.input, strings.Contains); err != nil {
				t.Fatalf("ParseFrom(%q) unexpected error: %v", tc.input, err)
			}

			if d.String() != tc.wantCanonical {
				t.Errorf("ParseFrom(%q): String() = %q, want %q", tc.input, d.String(), tc.wantCanonical)
			}
		})
	}
}
