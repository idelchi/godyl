package checksum_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/goccy/go-yaml"
	"github.com/google/go-cmp/cmp"

	"github.com/idelchi/godyl/internal/tools/checksum"
)

func TestIsChecksumLike(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input checksum.Checksums
		want  checksum.Checksums
	}{
		{
			name:  "checksums.txt",
			input: checksum.Checksums{"tool.tar.gz", "checksums.txt"},
			want:  checksum.Checksums{"checksums.txt"},
		},
		{
			name:  "tool.sha256",
			input: checksum.Checksums{"tool.tar.gz", "tool.sha256"},
			want:  checksum.Checksums{"tool.sha256"},
		},
		{
			name:  "SHA256SUMS",
			input: checksum.Checksums{"tool.tar.gz", "SHA256SUMS"},
			want:  checksum.Checksums{"SHA256SUMS"},
		},
		{
			name:  "no checksum files",
			input: checksum.Checksums{"tool.tar.gz", "README.md"},
			want:  nil,
		},
		{
			name:  "multiple",
			input: checksum.Checksums{"checksums.txt", "SHA256SUMS", "tool.tar.gz"},
			want:  checksum.Checksums{"checksums.txt", "SHA256SUMS"},
		},
		{
			name:  "empty",
			input: checksum.Checksums{},
			want:  nil,
		},
		{
			name:  "md5sums",
			input: checksum.Checksums{"tool.tar.gz", "md5sums.txt"},
			want:  checksum.Checksums{"md5sums.txt"},
		},
		{
			name:  "digest file",
			input: checksum.Checksums{"tool.tar.gz", "tool.digest"},
			want:  checksum.Checksums{"tool.digest"},
		},
		{
			// IsChecksumLike uses strings.ToLower before matching indicators,
			// so mixed-case extensions like ".SHA256" are still detected.
			name:  "mixed-case mid-extension",
			input: checksum.Checksums{"tool.tar.gz", "Tool.SHA256"},
			want:  checksum.Checksums{"Tool.SHA256"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := tc.input.IsChecksumLike()
			if !slices.Equal([]string(got), []string(tc.want)) {
				t.Errorf("IsChecksumLike() mismatch (-want +got):\n%s", cmp.Diff(tc.want, got))
			}
		})
	}
}

// TestPreferredSingleNonScoring documents the early-return behavior of
// Preferred() when the Checksums slice contains exactly one element: the
// element is returned unconditionally, even if it is not a checksum-like file
// (e.g. an archive). This is the len(cs)==1 fast-path in the implementation.
func TestPreferredSingleNonScoring(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		cs    checksum.Checksums
		asset string
		want  string
	}{
		{
			// The only element is a plain archive, not a checksum file. The
			// fast-path returns it directly without scoring.
			name:  "single non-checksum file is returned unconditionally",
			cs:    checksum.Checksums{"archive.tar.gz"},
			asset: "tool",
			want:  "archive.tar.gz",
		},
		{
			// Even a binary with no checksum indicators is returned when it is
			// the only entry.
			name:  "single binary asset returned unconditionally",
			cs:    checksum.Checksums{"mytool"},
			asset: "other-asset",
			want:  "mytool",
		},
		{
			// A real checksum file is also covered by the early-return path.
			name:  "single checksum file returned via early-return path",
			cs:    checksum.Checksums{"SHA256SUMS"},
			asset: "tool.tar.gz",
			want:  "SHA256SUMS",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := tc.cs.Preferred(tc.asset)
			if got != tc.want {
				t.Errorf("Preferred(%q) = %q, want %q", tc.asset, got, tc.want)
			}
		})
	}
}

func TestIndicators(t *testing.T) {
	t.Parallel()

	got := checksum.Indicators()

	if len(got) == 0 {
		t.Fatal("Indicators() returned empty slice")
	}

	// Spot-check that well-known indicators are present.
	want := []string{"checksum", "sha256", "md5", "digest", "sums"}
	for _, w := range want {
		if !slices.Contains(got, w) {
			t.Errorf("Indicators() missing %q", w)
		}
	}
}

func TestChecksumTypeString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input checksum.Type
		want  string
	}{
		{checksum.SHA256, "sha256"},
		{checksum.SHA512, "sha512"},
		{checksum.SHA1, "sha1"},
		{checksum.MD5, "md5"},
		{checksum.File, "file"},
		{checksum.None, "none"},
	}

	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			t.Parallel()

			if got := tc.input.String(); got != tc.want {
				t.Errorf("Type(%q).String() = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestChecksumIsSet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{name: "non-empty value is set", value: "abc123", want: true},
		{name: "empty value is not set", value: "", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c := checksum.Checksum{Value: tc.value}
			if got := c.IsSet(); got != tc.want {
				t.Errorf("IsSet() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestChecksumIsMandatory(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		typ  checksum.Type
		want bool
	}{
		{name: "sha256 is mandatory", typ: checksum.SHA256, want: true},
		{name: "file is mandatory", typ: checksum.File, want: true},
		{name: "none is not mandatory", typ: checksum.None, want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c := checksum.Checksum{Type: tc.typ}
			if got := c.IsMandatory(); got != tc.want {
				t.Errorf("IsMandatory() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestChecksumToQuery(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		typ   checksum.Type
		value string
		want  string
	}{
		{
			name:  "sha256 query",
			typ:   checksum.SHA256,
			value: "abc123",
			want:  "checksum=sha256:abc123",
		},
		{
			name:  "md5 query",
			typ:   checksum.MD5,
			value: "deadbeef",
			want:  "checksum=md5:deadbeef",
		},
		{
			name:  "empty value produces trailing colon",
			typ:   checksum.SHA256,
			value: "",
			want:  "checksum=sha256:",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c := checksum.Checksum{Type: tc.typ, Value: tc.value}
			if got := c.ToQuery(); got != tc.want {
				t.Errorf("ToQuery() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestPreferred(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		cs    checksum.Checksums
		asset string
		want  string
	}{
		{
			name:  "single entry",
			cs:    checksum.Checksums{"checksums.txt"},
			asset: "tool.tar.gz",
			want:  "checksums.txt",
		},
		{
			name:  "empty",
			cs:    checksum.Checksums{},
			asset: "tool.tar.gz",
			want:  "",
		},
		{
			name:  "prefers name match",
			cs:    checksum.Checksums{"checksums.txt", "tool_checksums.txt"},
			asset: "tool",
			want:  "tool_checksums.txt",
		},
		{
			name:  "generic only",
			cs:    checksum.Checksums{"checksums.txt", "sums.txt"},
			asset: "unrelated.tar.gz",
			want:  "checksums.txt",
		},
		{
			name:  "none scores above zero",
			cs:    checksum.Checksums{"archive.tar.gz", "binary.zip"},
			asset: "tool",
			want:  "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := tc.cs.Preferred(tc.asset)
			if got != tc.want {
				t.Errorf("Preferred(%q) = %q, want %q", tc.asset, got, tc.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Checksum.Resolve tests
// ---------------------------------------------------------------------------

func TestChecksumResolve(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     checksum.Checksum
		wantType  checksum.Type
		wantValue string
		wantErr   string
	}{
		{
			name:     "none type is a no-op",
			input:    checksum.Checksum{Type: checksum.None, Value: "anything"},
			wantType: checksum.None,
		},
		{
			name:     "file type without entry is a no-op",
			input:    checksum.Checksum{Type: checksum.File, Value: ""},
			wantType: checksum.File,
		},
		{
			name:    "file type with entry returns error",
			input:   checksum.Checksum{Type: checksum.File, Entry: "tool.tar.gz"},
			wantErr: "cannot use 'entry'",
		},
		{
			name:      "sha256: prefix strips algo and sets type",
			input:     checksum.Checksum{Type: checksum.SHA256, Value: "sha256:abc123"},
			wantType:  checksum.SHA256,
			wantValue: "abc123",
		},
		{
			name:      "sha512: prefix strips algo and sets type",
			input:     checksum.Checksum{Type: checksum.SHA256, Value: "sha512:def456"},
			wantType:  checksum.SHA512,
			wantValue: "def456",
		},
		{
			name:      "md5: prefix strips algo and sets type",
			input:     checksum.Checksum{Type: checksum.SHA256, Value: "md5:aabbcc"},
			wantType:  checksum.MD5,
			wantValue: "aabbcc",
		},
		{
			name:      "sha1: prefix strips algo and sets type",
			input:     checksum.Checksum{Type: checksum.SHA256, Value: "sha1:112233"},
			wantType:  checksum.SHA1,
			wantValue: "112233",
		},
		{
			name:      "sha256: prefix with whitespace trims value",
			input:     checksum.Checksum{Type: checksum.SHA256, Value: "sha256:  abc123  "},
			wantType:  checksum.SHA256,
			wantValue: "abc123",
		},
		{
			name:      "plain value without prefix is unchanged",
			input:     checksum.Checksum{Type: checksum.SHA256, Value: "plainvalue"},
			wantType:  checksum.SHA256,
			wantValue: "plainvalue",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c := tc.input
			err := c.Resolve(false)

			if tc.wantErr != "" {
				if err == nil {
					t.Fatalf("Resolve() expected error containing %q, got nil", tc.wantErr)
				}

				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Errorf("Resolve() error = %q, want containing %q", err.Error(), tc.wantErr)
				}

				return
			}

			if err != nil {
				t.Fatalf("Resolve() unexpected error: %v", err)
			}

			if c.Type != tc.wantType {
				t.Errorf("Resolve() Type = %q, want %q", c.Type, tc.wantType)
			}

			if tc.wantValue != "" && c.Value != tc.wantValue {
				t.Errorf("Resolve() Value = %q, want %q", c.Value, tc.wantValue)
			}
		})
	}
}

func TestChecksumResolvePath(t *testing.T) {
	t.Parallel()

	t.Run("reads checksum from file by entry", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		checksumFile := filepath.Join(dir, "checksums.txt")

		// Use a 64-char hex string so ParseChecksumFile recognises the GNU format.
		hash := strings.Repeat("ab", 32) // 64 hex chars

		if err := os.WriteFile(checksumFile, []byte(hash+"  tool.tar.gz\n"), 0o644); err != nil {
			t.Fatalf("WriteFile() error: %v", err)
		}

		c := checksum.Checksum{Type: checksum.SHA256, Value: "path:" + checksumFile, Entry: "tool.tar.gz"}

		if err := c.Resolve(false); err != nil {
			t.Fatalf("Resolve() unexpected error: %v", err)
		}

		if c.Value != hash {
			t.Errorf("Resolve() Value = %q, want %q", c.Value, hash)
		}
	})

	t.Run("reads raw checksum from file without entry", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		checksumFile := filepath.Join(dir, "sha.txt")

		if err := os.WriteFile(checksumFile, []byte("deadbeef\n"), 0o644); err != nil {
			t.Fatalf("WriteFile() error: %v", err)
		}

		c := checksum.Checksum{Type: checksum.SHA256, Value: "path:" + checksumFile}

		if err := c.Resolve(false); err != nil {
			t.Fatalf("Resolve() unexpected error: %v", err)
		}

		if c.Value != "deadbeef" {
			t.Errorf("Resolve() Value = %q, want %q", c.Value, "deadbeef")
		}
	})

	t.Run("entry not found in file returns error", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		checksumFile := filepath.Join(dir, "checksums.txt")

		if err := os.WriteFile(checksumFile, []byte("abc123  other.tar.gz\n"), 0o644); err != nil {
			t.Fatalf("WriteFile() error: %v", err)
		}

		c := checksum.Checksum{Type: checksum.SHA256, Value: "path:" + checksumFile, Entry: "missing.tar.gz"}

		if err := c.Resolve(false); err == nil {
			t.Fatal("Resolve() expected error for missing entry, got nil")
		}
	})

	t.Run("nonexistent path returns error", func(t *testing.T) {
		t.Parallel()

		c := checksum.Checksum{Type: checksum.SHA256, Value: "path:/nonexistent/file.txt"}

		if err := c.Resolve(false); err == nil {
			t.Fatal("Resolve() expected error for nonexistent path, got nil")
		}
	})
}

func TestChecksumResolveURL(t *testing.T) {
	t.Parallel()

	const hash = "abc123abc123abc123abc123abc123abc123abc123abc123abc123abc123abcd"

	t.Run("downloads and extracts checksum by entry", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(hash + "  tool-linux-amd64.tar.gz\n"))
		}))
		t.Cleanup(srv.Close)

		c := checksum.Checksum{
			Type:  checksum.SHA256,
			Value: "url:" + srv.URL + "/checksums.txt",
			Entry: "tool-linux-amd64.tar.gz",
		}

		if err := c.Resolve(false); err != nil {
			t.Fatalf("Resolve() unexpected error: %v", err)
		}

		if c.Value != hash {
			t.Errorf("Resolve() Value = %q, want %q", c.Value, hash)
		}
	})

	t.Run("downloads raw checksum without entry", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(hash + "\n"))
		}))
		t.Cleanup(srv.Close)

		c := checksum.Checksum{
			Type:  checksum.SHA256,
			Value: "url:" + srv.URL + "/sha.txt",
		}

		if err := c.Resolve(false); err != nil {
			t.Fatalf("Resolve() unexpected error: %v", err)
		}

		if c.Value != hash {
			t.Errorf("Resolve() Value = %q, want %q", c.Value, hash)
		}
	})

	t.Run("entry not found in URL response returns error", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(hash + "  other-tool.tar.gz\n"))
		}))
		t.Cleanup(srv.Close)

		c := checksum.Checksum{
			Type:  checksum.SHA256,
			Value: "url:" + srv.URL + "/checksums.txt",
			Entry: "missing.tar.gz",
		}

		if err := c.Resolve(false); err == nil {
			t.Fatal("Resolve() expected error for missing entry, got nil")
		}
	})
}

// ---------------------------------------------------------------------------
// Checksum UnmarshalYAML tests
// ---------------------------------------------------------------------------

func TestChecksumUnmarshalYAML(t *testing.T) {
	t.Parallel()

	t.Run("scalar string populates Type via single tag", func(t *testing.T) {
		t.Parallel()

		var c checksum.Checksum
		if err := yaml.Unmarshal([]byte(`sha256`), &c); err != nil {
			t.Fatalf("Unmarshal() unexpected error: %v", err)
		}

		if c.Type != checksum.SHA256 {
			t.Errorf("Unmarshal(\"sha256\").Type = %q, want %q", c.Type, checksum.SHA256)
		}
	})

	t.Run("map form populates all fields", func(t *testing.T) {
		t.Parallel()

		input := heredoc.Doc(`
			type: sha256
			value: abc123
			pattern: "checksums*"
		`)

		var c checksum.Checksum
		if err := yaml.Unmarshal([]byte(input), &c); err != nil {
			t.Fatalf("Unmarshal() unexpected error: %v", err)
		}

		if c.Type != checksum.SHA256 {
			t.Errorf("Type = %q, want %q", c.Type, checksum.SHA256)
		}

		if c.Value != "abc123" {
			t.Errorf("Value = %q, want %q", c.Value, "abc123")
		}

		if c.Pattern != "checksums*" {
			t.Errorf("Pattern = %q, want %q", c.Pattern, "checksums*")
		}
	})

	t.Run("none type unmarshals correctly", func(t *testing.T) {
		t.Parallel()

		var c checksum.Checksum
		if err := yaml.Unmarshal([]byte(`none`), &c); err != nil {
			t.Fatalf("Unmarshal() unexpected error: %v", err)
		}

		if c.Type != checksum.None {
			t.Errorf("Type = %q, want %q", c.Type, checksum.None)
		}
	})
}
