package checksum_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/idelchi/godyl/internal/tools/checksum"
)

func TestParseChecksumFile(t *testing.T) {
	t.Parallel()

	const (
		hash1 = "abc123abc123abc123abc123abc123abc123abc123abc123abc123abc123abcd" // sha256 hex (64 chars)
		hash2 = "def456def456def456def456def456def456def456def456def456def456deff" // sha256 hex (64 chars)
	)

	tests := []struct {
		name  string
		input string
		want  map[string]string
	}{
		{
			name:  "GNU format single entry",
			input: hash1 + "  filename.tar.gz",
			want:  map[string]string{"filename.tar.gz": hash1},
		},
		{
			name:  "GNU format binary mode indicator stripped",
			input: hash1 + " *filename.tar.gz",
			want:  map[string]string{"filename.tar.gz": hash1},
		},
		{
			name: "GNU format multi-line two entries",
			input: strings.Join([]string{
				hash1 + "  file-one.tar.gz",
				hash2 + "  file-two.tar.gz",
			}, "\n"),
			want: map[string]string{
				"file-one.tar.gz": hash1,
				"file-two.tar.gz": hash2,
			},
		},
		{
			name:  "BSD format single entry",
			input: "SHA256 (file.tar.gz) = " + hash1,
			want:  map[string]string{"file.tar.gz": hash1},
		},
		{
			name: "BSD format multi-line two entries",
			input: strings.Join([]string{
				"SHA256 (file-one.tar.gz) = " + hash1,
				"SHA256 (file-two.tar.gz) = " + hash2,
			}, "\n"),
			want: map[string]string{
				"file-one.tar.gz": hash1,
				"file-two.tar.gz": hash2,
			},
		},
		{
			// ParseChecksumFile detects format from the first non-junk line.
			// A file with BSD entries followed by GNU entries is dispatched as
			// BSD-format, so only the BSD lines are parsed.
			name: "mixed BSD then GNU lines parsed as BSD format",
			input: strings.Join([]string{
				"SHA256 (file-one.tar.gz) = " + hash1,
				hash2 + "  file-two.tar.gz",
			}, "\n"),
			want: map[string]string{
				"file-one.tar.gz": hash1,
				// file-two.tar.gz is skipped: the BSD regex doesn't match GNU lines
			},
		},
		{
			name:  "empty input returns empty map",
			input: "",
			want:  map[string]string{},
		},
		{
			name:  "whitespace only returns empty map",
			input: "   \n\t\n   ",
			want:  map[string]string{},
		},
		{
			name: "junk first line then valid GNU entry",
			input: strings.Join([]string{
				"this is not a checksum line",
				hash1 + "  real-file.tar.gz",
			}, "\n"),
			want: map[string]string{"real-file.tar.gz": hash1},
		},
		{
			// A trailing newline after a valid GNU entry must not drop the entry.
			name:  "trailing newline does not lose entry",
			input: hash1 + "  filename.tar.gz\n",
			want:  map[string]string{"filename.tar.gz": hash1},
		},
		{
			// All lines are junk (no BSD or GNU pattern matches), so the result is an empty map.
			name: "all junk lines returns empty map",
			input: strings.Join([]string{
				"this is not a checksum",
				"neither is this",
				"or this one",
			}, "\n"),
			want: map[string]string{},
		},
		{
			// GNU binary mode uses "* " prefix; some tools emit "**" (double star).
			// The regex [* ] matches the first '*', and .+ captures "*filename.tar.gz".
			// TrimPrefix strips one leading '*', yielding "filename.tar.gz".
			name:  "GNU format name with double star strips one star",
			input: hash1 + " **filename.tar.gz",
			want:  map[string]string{"filename.tar.gz": hash1},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := checksum.ParseChecksumFile(tc.input)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("ParseChecksumFile(%q) mismatch (-want +got):\n%s", tc.input, diff)
			}
		})
	}
}

func TestInferAlgoFromHex(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "64-char hex infers sha256",
			input: strings.Repeat("a", 64),
			want:  "sha256",
		},
		{
			name:  "128-char hex infers sha512",
			input: strings.Repeat("a", 128),
			want:  "sha512",
		},
		{
			name:  "40-char hex infers sha1",
			input: strings.Repeat("a", 40),
			want:  "sha1",
		},
		{
			name:  "32-char hex infers md5",
			input: strings.Repeat("a", 32),
			want:  "md5",
		},
		{
			name:  "10-char hex falls through to default sha256",
			input: strings.Repeat("a", 10),
			want:  "sha256",
		},
		{
			name:  "empty string falls through to default sha256",
			input: "",
			want:  "sha256",
		},
		{
			// InferAlgoFromHex dispatches on length only; it does not validate
			// whether the characters are valid hex digits. A 64-character string
			// of non-hex characters still returns sha256.
			name:  "64-char non-hex string still infers sha256 (length-only behavior)",
			input: strings.Repeat("z", 64),
			want:  "sha256",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := checksum.InferAlgoFromHex(tc.input)
			if got != tc.want {
				t.Errorf("InferAlgoFromHex(%q): got %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
