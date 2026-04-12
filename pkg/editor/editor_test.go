package editor_test

import (
	"path/filepath"
	"testing"

	"github.com/idelchi/godyl/pkg/editor"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/unmarshal"
)

const wantKeyValue = "key: value\n"

func TestMerge(t *testing.T) {
	t.Parallel()

	t.Run("into empty file", func(t *testing.T) {
		t.Parallel()

		f := file.New(filepath.Join(t.TempDir(), "empty.yml"))

		if err := f.Create(); err != nil {
			t.Fatalf("Create(): unexpected error: %v", err)
		}

		if err := editor.New(f).Merge(map[string]any{"key": "value"}); err != nil {
			t.Fatalf("Merge(): unexpected error: %v", err)
		}

		got, err := f.ReadString()
		if err != nil {
			t.Fatalf("ReadString(): unexpected error: %v", err)
		}

		if got != wantKeyValue {
			t.Errorf("Merge() wrote %q, want %q", got, wantKeyValue)
		}
	})

	t.Run("into nonexistent file", func(t *testing.T) {
		t.Parallel()

		f := file.New(filepath.Join(t.TempDir(), "subdir", "new.yml"))

		if err := editor.New(f).Merge(map[string]any{"key": "value"}); err != nil {
			t.Fatalf("Merge(): unexpected error: %v", err)
		}

		got, err := f.ReadString()
		if err != nil {
			t.Fatalf("ReadString(): unexpected error: %v", err)
		}

		if got != wantKeyValue {
			t.Errorf("Merge() wrote %q, want %q", got, wantKeyValue)
		}
	})

	t.Run("into existing data", func(t *testing.T) {
		t.Parallel()

		f := file.New(filepath.Join(t.TempDir(), "existing.yml"))

		if err := f.Write([]byte("existing: 'data'\n")); err != nil {
			t.Fatalf("Write(): unexpected error: %v", err)
		}

		if err := editor.New(f).Merge(map[string]any{"new": "entry"}); err != nil {
			t.Fatalf("Merge(): unexpected error: %v", err)
		}

		data, err := f.Read()
		if err != nil {
			t.Fatalf("Read(): unexpected error: %v", err)
		}

		got := make(map[string]any)
		if err := unmarshal.Lax(data, &got); err != nil {
			t.Fatalf("unmarshal.Lax(): unexpected error: %v", err)
		}

		if got["existing"] != "data" {
			t.Errorf("got[existing] = %v, want %q", got["existing"], "data")
		}

		if got["new"] != "entry" {
			t.Errorf("got[new] = %v, want %q", got["new"], "entry")
		}
	})
}

func TestWrite_EmptyFile(t *testing.T) {
	t.Parallel()

	f := file.New(filepath.Join(t.TempDir(), "empty.yml"))

	if err := f.Create(); err != nil {
		t.Fatalf("Create(): unexpected error: %v", err)
	}

	if err := editor.New(f).Write(map[string]any{"key": "value"}); err != nil {
		t.Fatalf("Write(): unexpected error: %v", err)
	}

	got, err := f.ReadString()
	if err != nil {
		t.Fatalf("ReadString(): unexpected error: %v", err)
	}

	if got != wantKeyValue {
		t.Errorf("Write() wrote %q, want %q", got, wantKeyValue)
	}
}
