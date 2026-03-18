package iutils_test

import (
	"errors"
	"path/filepath"
	"slices"
	"testing"

	"github.com/idelchi/godyl/internal/iutils"
	"github.com/idelchi/godyl/pkg/path/file"
)

func TestPrefix(t *testing.T) {
	t.Parallel()

	t.Run("Lower", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name  string
			input iutils.Prefix
			want  iutils.Prefix
		}{
			{name: "upper to lower", input: "GODYL_INSTALL", want: "godyl_install"},
			{name: "already lower", input: "godyl", want: "godyl"},
			{name: "empty", input: "", want: ""},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				got := tc.input.Lower()
				if got != tc.want {
					t.Errorf("Prefix(%q).Lower() = %q, want %q", tc.input, got, tc.want)
				}
			})
		}
	})

	t.Run("Upper", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name  string
			input iutils.Prefix
			want  iutils.Prefix
		}{
			{name: "lower to upper", input: "godyl", want: "GODYL"},
			{name: "already upper", input: "GODYL", want: "GODYL"},
			{name: "empty", input: "", want: ""},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				got := tc.input.Upper()
				if got != tc.want {
					t.Errorf("Prefix(%q).Upper() = %q, want %q", tc.input, got, tc.want)
				}
			})
		}
	})

	t.Run("RemovePrefix", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name   string
			input  iutils.Prefix
			prefix string
			want   iutils.Prefix
		}{
			{name: "removes matching prefix", input: "GODYL_INSTALL", prefix: "GODYL_", want: "INSTALL"},
			{name: "no-op when prefix absent", input: "GODYL_INSTALL", prefix: "MISSING_", want: "GODYL_INSTALL"},
			{name: "empty input", input: "", prefix: "GODYL_", want: ""},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				got := tc.input.RemovePrefix(tc.prefix)
				if got != tc.want {
					t.Errorf("Prefix(%q).RemovePrefix(%q) = %q, want %q", tc.input, tc.prefix, got, tc.want)
				}
			})
		}
	})

	t.Run("RemoveSuffix", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name   string
			input  iutils.Prefix
			suffix string
			want   iutils.Prefix
		}{
			{name: "removes matching suffix", input: "GODYL_INSTALL", suffix: "_INSTALL", want: "GODYL"},
			{name: "no-op when suffix absent", input: "GODYL_INSTALL", suffix: "_MISSING", want: "GODYL_INSTALL"},
			{name: "empty input", input: "", suffix: "_INSTALL", want: ""},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				got := tc.input.RemoveSuffix(tc.suffix)
				if got != tc.want {
					t.Errorf("Prefix(%q).RemoveSuffix(%q) = %q, want %q", tc.input, tc.suffix, got, tc.want)
				}
			})
		}
	})

	t.Run("WithUnderscores", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name  string
			input iutils.Prefix
			want  iutils.Prefix
		}{
			{name: "dot replaced with underscore", input: "godyl.install", want: "godyl_install"},
			{name: "already underscored", input: "godyl_install", want: "godyl_install"},
			{name: "empty", input: "", want: ""},
			{name: "multiple dots replaced", input: "godyl.install.thing", want: "godyl_install_thing"},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				got := tc.input.WithUnderscores()
				if got != tc.want {
					t.Errorf("Prefix(%q).WithUnderscores() = %q, want %q", tc.input, got, tc.want)
				}
			})
		}
	})

	t.Run("Scoped", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name  string
			input iutils.Prefix
			want  iutils.Prefix
		}{
			{name: "upper appends underscore", input: "GODYL", want: "GODYL_"},
			{name: "lower appends underscore", input: "godyl", want: "godyl_"},
			{name: "empty appends underscore", input: "", want: "_"},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				got := tc.input.Scoped()
				if got != tc.want {
					t.Errorf("Prefix(%q).Scoped() = %q, want %q", tc.input, got, tc.want)
				}
			})
		}
	})
}

func TestMatchEnvToFlag(t *testing.T) {
	t.Parallel()

	t.Run("GODYL_ prefix", func(t *testing.T) {
		t.Parallel()

		matcher := iutils.MatchEnvToFlag("GODYL_")

		tests := []struct {
			name  string
			input string
			want  string
		}{
			{name: "multi-segment env var", input: "GODYL_GITHUB_TOKEN", want: "github-token"},
			{name: "two-segment env var", input: "GODYL_NO_CACHE", want: "no-cache"},
			{name: "log level", input: "GODYL_LOG_LEVEL", want: "log-level"},
			{name: "single segment", input: "GODYL_PARALLEL", want: "parallel"},
			{name: "env var exactly equal to prefix yields empty string", input: "GODYL_", want: ""},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				got := matcher(tc.input)
				if got != tc.want {
					t.Errorf("MatchEnvToFlag(\"GODYL_\")(%q) = %q, want %q", tc.input, got, tc.want)
				}
			})
		}
	})

	t.Run("non-matching prefix returns full lowercased hyphenated string", func(t *testing.T) {
		t.Parallel()

		matcher := iutils.MatchEnvToFlag("GODYL_")

		// "OTHER_LOG_LEVEL" does not start with "godyl_" after lowercasing.
		// TrimPrefix is a no-op, so the result is the full lowercased+hyphenated string.
		got := matcher("OTHER_LOG_LEVEL")
		want := "other-log-level"

		if got != want {
			t.Errorf("MatchEnvToFlag(\"GODYL_\")(\"OTHER_LOG_LEVEL\") = %q, want %q", got, want)
		}
	})

	t.Run("mixed-case prefix is normalised", func(t *testing.T) {
		t.Parallel()

		matcher := iutils.MatchEnvToFlag("MyApp_")

		tests := []struct {
			name  string
			input string
			want  string
		}{
			{name: "uppercase input with mixed-case prefix", input: "MYAPP_LOG_LEVEL", want: "log-level"},
			{name: "matching case input", input: "myapp_token", want: "token"},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				got := matcher(tc.input)
				if got != tc.want {
					t.Errorf("MatchEnvToFlag(\"MyApp_\")(%q) = %q, want %q", tc.input, got, tc.want)
				}
			})
		}
	})
}

func TestSplitTags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       []string
		wantInclude []string
		wantExclude []string
	}{
		{
			name:        "plain tags go to include",
			input:       []string{"linux", "amd64"},
			wantInclude: []string{"linux", "amd64"},
		},
		{
			name:        "bang-prefixed tags go to exclude",
			input:       []string{"!windows", "!arm"},
			wantExclude: []string{"windows", "arm"},
		},
		{
			name:        "mixed tags are split correctly",
			input:       []string{"linux", "!windows", "amd64", "!arm"},
			wantInclude: []string{"linux", "amd64"},
			wantExclude: []string{"windows", "arm"},
		},
		{
			name:  "empty input produces empty result",
			input: []string{},
		},
		{
			name:  "nil input produces empty result",
			input: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := iutils.SplitTags(tc.input)

			if !slices.Equal([]string(got.Include), tc.wantInclude) {
				t.Errorf("SplitTags(%v).Include = %v, want %v", tc.input, got.Include, tc.wantInclude)
			}

			if !slices.Equal([]string(got.Exclude), tc.wantExclude) {
				t.Errorf("SplitTags(%v).Exclude = %v, want %v", tc.input, got.Exclude, tc.wantExclude)
			}
		})
	}
}

func TestAny(t *testing.T) {
	t.Parallel()

	t.Run("returns first non-zero string", func(t *testing.T) {
		t.Parallel()

		got := iutils.Any("", "", "hello", "world")
		if got != "hello" {
			t.Errorf("Any() = %q, want %q", got, "hello")
		}
	})

	t.Run("returns first non-zero int", func(t *testing.T) {
		t.Parallel()

		got := iutils.Any(0, 0, 42, 99)
		if got != 42 {
			t.Errorf("Any() = %d, want %d", got, 42)
		}
	})

	t.Run("all zero returns zero", func(t *testing.T) {
		t.Parallel()

		got := iutils.Any("", "", "")
		if got != "" {
			t.Errorf("Any() = %q, want empty string", got)
		}
	})

	t.Run("no args returns zero", func(t *testing.T) {
		t.Parallel()

		got := iutils.Any[int]()
		if got != 0 {
			t.Errorf("Any() = %d, want 0", got)
		}
	})

	t.Run("first arg is non-zero", func(t *testing.T) {
		t.Parallel()

		got := iutils.Any("first", "second")
		if got != "first" {
			t.Errorf("Any() = %q, want %q", got, "first")
		}
	})
}

func TestBytesSource(t *testing.T) {
	t.Parallel()

	t.Run("returns stored data", func(t *testing.T) {
		t.Parallel()

		data := []byte("hello world")
		src := iutils.BytesSource{Data: data}

		got, err := src.Read()
		if err != nil {
			t.Fatalf("Read() unexpected error: %v", err)
		}

		if string(got) != string(data) {
			t.Errorf("Read() = %q, want %q", got, data)
		}
	})

	t.Run("nil data returns nil", func(t *testing.T) {
		t.Parallel()

		src := iutils.BytesSource{}

		got, err := src.Read()
		if err != nil {
			t.Fatalf("Read() unexpected error: %v", err)
		}

		if got != nil {
			t.Errorf("Read() = %v, want nil", got)
		}
	})
}

func TestMultiSource(t *testing.T) {
	t.Parallel()

	t.Run("concatenates sources with newlines", func(t *testing.T) {
		t.Parallel()

		ms := iutils.NewMultiSource(
			iutils.BytesSource{Data: []byte("aaa")},
			iutils.BytesSource{Data: []byte("bbb")},
			iutils.BytesSource{Data: []byte("ccc")},
		)

		got, err := ms.Read()
		if err != nil {
			t.Fatalf("Read() unexpected error: %v", err)
		}

		want := "aaa\nbbb\nccc"
		if string(got) != want {
			t.Errorf("Read() = %q, want %q", got, want)
		}
	})

	t.Run("single source has no extra newline", func(t *testing.T) {
		t.Parallel()

		ms := iutils.NewMultiSource(
			iutils.BytesSource{Data: []byte("only")},
		)

		got, err := ms.Read()
		if err != nil {
			t.Fatalf("Read() unexpected error: %v", err)
		}

		if string(got) != "only" {
			t.Errorf("Read() = %q, want %q", got, "only")
		}
	})

	t.Run("empty sources returns empty", func(t *testing.T) {
		t.Parallel()

		ms := iutils.NewMultiSource()

		got, err := ms.Read()
		if err != nil {
			t.Fatalf("Read() unexpected error: %v", err)
		}

		if len(got) != 0 {
			t.Errorf("Read() = %q, want empty", got)
		}
	})
}

func TestFileSourceRead(t *testing.T) {
	t.Parallel()

	t.Run("reads file contents", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		f := file.New(filepath.Join(dir, "input.yaml"))

		if err := f.Write([]byte("key: value")); err != nil {
			t.Fatalf("Write() error: %v", err)
		}

		src := iutils.FileSource{File: f}

		got, err := src.Read()
		if err != nil {
			t.Fatalf("Read() unexpected error: %v", err)
		}

		if string(got) != "key: value" {
			t.Errorf("Read() = %q, want %q", got, "key: value")
		}
	})

	t.Run("nonexistent file returns error", func(t *testing.T) {
		t.Parallel()

		src := iutils.FileSource{File: file.New("/nonexistent/path.txt")}

		_, err := src.Read()
		if err == nil {
			t.Fatal("Read() expected error for nonexistent file, got nil")
		}
	})
}

func TestGetSourceFromPath(t *testing.T) {
	t.Parallel()

	t.Run("file path returns FileSource", func(t *testing.T) {
		t.Parallel()

		src, err := iutils.GetSourceFromPath("/some/path.yaml")
		if err != nil {
			t.Fatalf("GetSourceFromPath() unexpected error: %v", err)
		}

		if _, ok := src.(iutils.FileSource); !ok {
			t.Errorf("GetSourceFromPath(\"/some/path.yaml\") = %T, want FileSource", src)
		}
	})

	t.Run("dash without piped stdin returns error", func(t *testing.T) {
		t.Parallel()

		// In a test context stdin is not piped, so "-" should error.
		_, err := iutils.GetSourceFromPath("-")
		if err == nil {
			t.Fatal("GetSourceFromPath(\"-\") expected error when stdin is not piped, got nil")
		}

		if !errors.Is(err, iutils.ErrInvalidSource) {
			t.Errorf("GetSourceFromPath(\"-\") error = %v, want wrapping ErrInvalidSource", err)
		}
	})
}

func TestReadPaths(t *testing.T) {
	t.Parallel()

	t.Run("reads single file", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		p := filepath.Join(dir, "a.yaml")

		if err := file.New(p).Write([]byte("hello")); err != nil {
			t.Fatalf("Write() error: %v", err)
		}

		got, err := iutils.ReadPaths(p)
		if err != nil {
			t.Fatalf("ReadPaths() unexpected error: %v", err)
		}

		if string(got) != "hello" {
			t.Errorf("ReadPaths() = %q, want %q", got, "hello")
		}
	})

	t.Run("concatenates multiple files", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		p1 := filepath.Join(dir, "a.yaml")
		p2 := filepath.Join(dir, "b.yaml")

		if err := file.New(p1).Write([]byte("aaa")); err != nil {
			t.Fatalf("Write() error: %v", err)
		}

		if err := file.New(p2).Write([]byte("bbb")); err != nil {
			t.Fatalf("Write() error: %v", err)
		}

		got, err := iutils.ReadPaths(p1, p2)
		if err != nil {
			t.Fatalf("ReadPaths() unexpected error: %v", err)
		}

		if string(got) != "aaa\nbbb" {
			t.Errorf("ReadPaths() = %q, want %q", got, "aaa\nbbb")
		}
	})

	t.Run("no paths returns error", func(t *testing.T) {
		t.Parallel()

		_, err := iutils.ReadPaths()
		if err == nil {
			t.Fatal("ReadPaths() expected error with no paths, got nil")
		}

		if !errors.Is(err, iutils.ErrInvalidSource) {
			t.Errorf("ReadPaths() error = %v, want wrapping ErrInvalidSource", err)
		}
	})
}

func TestReadPathsOrDefault(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	defaultFile := filepath.Join(dir, "default.yaml")
	argFile := filepath.Join(dir, "arg.yaml")

	if err := file.New(defaultFile).Write([]byte("default")); err != nil {
		t.Fatalf("Write() error: %v", err)
	}

	if err := file.New(argFile).Write([]byte("from-arg")); err != nil {
		t.Fatalf("Write() error: %v", err)
	}

	t.Run("no args uses default path", func(t *testing.T) {
		t.Parallel()

		got, err := iutils.ReadPathsOrDefault(defaultFile)
		if err != nil {
			t.Fatalf("ReadPathsOrDefault() unexpected error: %v", err)
		}

		if string(got) != "default" {
			t.Errorf("ReadPathsOrDefault() = %q, want %q", got, "default")
		}
	})

	t.Run("args override default", func(t *testing.T) {
		t.Parallel()

		got, err := iutils.ReadPathsOrDefault(defaultFile, argFile)
		if err != nil {
			t.Fatalf("ReadPathsOrDefault() unexpected error: %v", err)
		}

		if string(got) != "from-arg" {
			t.Errorf("ReadPathsOrDefault() = %q, want %q", got, "from-arg")
		}
	})
}

func TestMerge(t *testing.T) {
	t.Parallel()

	type config struct {
		Name    string
		Count   int
		Enabled bool
	}

	t.Run("merges src into dst with override", func(t *testing.T) {
		t.Parallel()

		dst := config{Name: "old", Count: 1}
		src := config{Name: "new", Count: 2, Enabled: true}

		if err := iutils.Merge(&dst, &src); err != nil {
			t.Fatalf("Merge() unexpected error: %v", err)
		}

		if dst.Name != "new" {
			t.Errorf("Name = %q, want %q", dst.Name, "new")
		}

		if dst.Count != 2 {
			t.Errorf("Count = %d, want %d", dst.Count, 2)
		}

		if !dst.Enabled {
			t.Error("Enabled = false, want true")
		}
	})

	t.Run("zero src fields do not override dst", func(t *testing.T) {
		t.Parallel()

		dst := config{Name: "keep", Count: 5}
		src := config{Name: "", Count: 0}

		if err := iutils.Merge(&dst, &src); err != nil {
			t.Fatalf("Merge() unexpected error: %v", err)
		}

		// WithOverride + WithoutDereference: zero values in src do NOT override dst.
		if dst.Name != "keep" {
			t.Errorf("Name = %q, want %q (zero src should not override)", dst.Name, "keep")
		}

		if dst.Count != 5 {
			t.Errorf("Count = %d, want %d (zero src should not override)", dst.Count, 5)
		}
	})
}
