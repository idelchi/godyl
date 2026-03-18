package install_test

import (
	"testing"

	"github.com/idelchi/godyl/internal/tools/sources/install"
)

func TestSplitName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		input      string
		wantFirst  string
		wantSecond string
		wantErr    bool
	}{
		{
			name:       "valid owner/repo",
			input:      "owner/repo",
			wantFirst:  "owner",
			wantSecond: "repo",
		},
		{
			name:    "no slash",
			input:   "noslash",
			wantErr: true,
		},
		{
			name:    "too many parts",
			input:   "a/b/c",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			// SplitName splits on "/" and requires exactly two parts, so
			// "owner/" produces ["owner", ""] — two parts, so no error.
			// The second component is the empty string.
			name:       "trailing slash produces empty second component",
			input:      "owner/",
			wantFirst:  "owner",
			wantSecond: "",
		},
		{
			// SplitName splits on "/" and requires exactly two parts, so
			// "/repo" produces ["", "repo"] — two parts, so no error.
			// The first component is the empty string.
			name:       "leading slash produces empty first component",
			input:      "/repo",
			wantFirst:  "",
			wantSecond: "repo",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			first, second, err := install.SplitName(tc.input)

			if tc.wantErr {
				if err == nil {
					t.Errorf("SplitName(%q) expected error, got nil", tc.input)
				}

				return
			}

			if err != nil {
				t.Fatalf("SplitName(%q) unexpected error: %v", tc.input, err)
			}

			if first != tc.wantFirst {
				t.Errorf("SplitName(%q) first = %q, want %q", tc.input, first, tc.wantFirst)
			}

			if second != tc.wantSecond {
				t.Errorf("SplitName(%q) second = %q, want %q", tc.input, second, tc.wantSecond)
			}
		})
	}
}

// TestSplitNameUnicode verifies that SplitName and CutName handle inputs
// containing unicode characters or embedded whitespace correctly.
// Both functions operate on byte-level string splitting, so unicode and
// whitespace in the owner or repository components pass through unchanged.
func TestSplitNameUnicode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		input        string
		splitFirst   string
		splitSecond  string
		splitWantErr bool
		cutFirst     string
		cutSecond    string
		cutWantErr   bool
	}{
		{
			// Components with embedded spaces are split on '/' only.
			name:        "spaces in both components",
			input:       "a b/c d",
			splitFirst:  "a b",
			splitSecond: "c d",
			cutFirst:    "a b",
			cutSecond:   "c d",
		},
		{
			// Unicode owner and repo names are passed through as-is.
			name:        "unicode owner and repo",
			input:       "ünïcödé/répo",
			splitFirst:  "ünïcödé",
			splitSecond: "répo",
			cutFirst:    "ünïcödé",
			cutSecond:   "répo",
		},
		{
			// A tab character in the owner is preserved.
			name:        "tab character in owner",
			input:       "ow\tner/repo",
			splitFirst:  "ow\tner",
			splitSecond: "repo",
			cutFirst:    "ow\tner",
			cutSecond:   "repo",
		},
		{
			// CJK characters are valid in component names.
			name:        "CJK characters",
			input:       "所有者/仓库",
			splitFirst:  "所有者",
			splitSecond: "仓库",
			cutFirst:    "所有者",
			cutSecond:   "仓库",
		},
		{
			// Three slash-separated parts: SplitName errors, CutName returns remainder.
			name:         "unicode three parts",
			input:        "α/β/γ",
			splitWantErr: true,
			cutFirst:     "α",
			cutSecond:    "β/γ",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// --- SplitName ---
			sFirst, sSecond, sErr := install.SplitName(tc.input)

			if tc.splitWantErr {
				if sErr == nil {
					t.Errorf("SplitName(%q) expected error, got nil", tc.input)
				}
			} else {
				if sErr != nil {
					t.Fatalf("SplitName(%q) unexpected error: %v", tc.input, sErr)
				}

				if sFirst != tc.splitFirst {
					t.Errorf("SplitName(%q) first = %q, want %q", tc.input, sFirst, tc.splitFirst)
				}

				if sSecond != tc.splitSecond {
					t.Errorf("SplitName(%q) second = %q, want %q", tc.input, sSecond, tc.splitSecond)
				}
			}

			// --- CutName ---
			cFirst, cSecond, cErr := install.CutName(tc.input)

			if tc.cutWantErr {
				if cErr == nil {
					t.Errorf("CutName(%q) expected error, got nil", tc.input)
				}
			} else {
				if cErr != nil {
					t.Fatalf("CutName(%q) unexpected error: %v", tc.input, cErr)
				}

				if cFirst != tc.cutFirst {
					t.Errorf("CutName(%q) first = %q, want %q", tc.input, cFirst, tc.cutFirst)
				}

				if cSecond != tc.cutSecond {
					t.Errorf("CutName(%q) second = %q, want %q", tc.input, cSecond, tc.cutSecond)
				}
			}
		})
	}
}

func TestCutName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		input      string
		wantFirst  string
		wantSecond string
		wantErr    bool
	}{
		{
			name:       "valid owner/repo",
			input:      "owner/repo",
			wantFirst:  "owner",
			wantSecond: "repo",
		},
		{
			name:    "no slash",
			input:   "noslash",
			wantErr: true,
		},
		{
			name:       "three parts keeps remainder",
			input:      "a/b/c",
			wantFirst:  "a",
			wantSecond: "b/c",
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			// strings.Cut("owner/", "/") → first="owner", second="", found=true → no error.
			name:       "trailing slash produces empty second component",
			input:      "owner/",
			wantFirst:  "owner",
			wantSecond: "",
		},
		{
			// strings.Cut("/repo", "/") → first="", second="repo", found=true → no error.
			name:       "leading slash produces empty first component",
			input:      "/repo",
			wantFirst:  "",
			wantSecond: "repo",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			first, second, err := install.CutName(tc.input)

			if tc.wantErr {
				if err == nil {
					t.Errorf("CutName(%q) expected error, got nil", tc.input)
				}

				return
			}

			if err != nil {
				t.Fatalf("CutName(%q) unexpected error: %v", tc.input, err)
			}

			if first != tc.wantFirst {
				t.Errorf("CutName(%q) first = %q, want %q", tc.input, first, tc.wantFirst)
			}

			if second != tc.wantSecond {
				t.Errorf("CutName(%q) second = %q, want %q", tc.input, second, tc.wantSecond)
			}
		})
	}
}

func TestMetadataGet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		m    *install.Metadata
		key  string
		want string
	}{
		{
			name: "existing key returns value",
			m:    &install.Metadata{"version": "1.0"},
			key:  "version",
			want: "1.0",
		},
		{
			name: "missing key returns empty",
			m:    &install.Metadata{"version": "1.0"},
			key:  "absent",
			want: "",
		},
		{
			name: "nil metadata returns empty",
			m:    nil,
			key:  "anything",
			want: "",
		},
		{
			name: "empty metadata returns empty",
			m:    &install.Metadata{},
			key:  "anything",
			want: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := tc.m.Get(tc.key)
			if got != tc.want {
				t.Errorf("Get(%q) = %q, want %q", tc.key, got, tc.want)
			}
		})
	}
}

func TestMetadataSet(t *testing.T) {
	t.Parallel()

	t.Run("set on initialized metadata", func(t *testing.T) {
		t.Parallel()

		m := install.Metadata{"existing": "value"}
		m.Set("new", "data")

		if got := m.Get("new"); got != "data" {
			t.Errorf("Get(\"new\") after Set = %q, want %q", got, "data")
		}

		if got := m.Get("existing"); got != "value" {
			t.Errorf("Get(\"existing\") = %q, want %q (should be preserved)", got, "value")
		}
	})

	t.Run("set on nil metadata initializes it", func(t *testing.T) {
		t.Parallel()

		var m install.Metadata
		m.Set("key", "val")

		if got := m.Get("key"); got != "val" {
			t.Errorf("Get(\"key\") after Set on nil = %q, want %q", got, "val")
		}
	})

	t.Run("set overwrites existing key", func(t *testing.T) {
		t.Parallel()

		m := install.Metadata{"key": "old"}
		m.Set("key", "new")

		if got := m.Get("key"); got != "new" {
			t.Errorf("Get(\"key\") after overwrite = %q, want %q", got, "new")
		}
	})
}
