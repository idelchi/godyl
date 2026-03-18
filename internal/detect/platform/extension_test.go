package platform_test

import (
	"strings"
	"testing"

	"github.com/idelchi/godyl/internal/detect/platform"
)

func TestExtensionParseFrom(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		osName    string
		wantValue string
	}{
		{
			name:      "windows yields .exe",
			osName:    "windows",
			wantValue: ".exe",
		},
		{
			name:      "linux yields empty string",
			osName:    "linux",
			wantValue: "",
		},
		{
			name:      "darwin yields empty string",
			osName:    "darwin",
			wantValue: "",
		},
		{
			name:      "freebsd yields empty string",
			osName:    "freebsd",
			wantValue: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var o platform.OS

			if err := o.ParseFrom(tc.osName, strings.EqualFold); err != nil {
				t.Fatalf("OS.ParseFrom(%q) unexpected error: %v", tc.osName, err)
			}

			var e platform.Extension

			e.ParseFrom(o)

			if string(e) != tc.wantValue {
				t.Errorf("ParseFrom(%q): value = %q, want %q", tc.osName, string(e), tc.wantValue)
			}

			if e.String() != tc.wantValue {
				t.Errorf("ParseFrom(%q): String() = %q, want %q", tc.osName, e.String(), tc.wantValue)
			}
		})
	}
}

func TestExtensionIsNil(t *testing.T) {
	t.Parallel()

	t.Run("zero-value extension is nil", func(t *testing.T) {
		t.Parallel()

		var e platform.Extension

		if !e.IsNil() {
			t.Error("zero-value Extension.IsNil() = false, want true")
		}
	})

	t.Run("windows extension is not nil", func(t *testing.T) {
		t.Parallel()

		var o platform.OS

		if err := o.ParseFrom("windows", strings.EqualFold); err != nil {
			t.Fatalf("OS.ParseFrom(%q) unexpected error: %v", "windows", err)
		}

		var e platform.Extension

		e.ParseFrom(o)

		if e.IsNil() {
			t.Error("windows Extension.IsNil() = true, want false")
		}
	})
}

func TestExtensionIsNilAllOS(t *testing.T) {
	t.Parallel()

	nonWindowsOSes := []string{"linux", "darwin", "freebsd"}

	for _, osName := range nonWindowsOSes {
		t.Run(osName+" extension is nil", func(t *testing.T) {
			t.Parallel()

			var o platform.OS

			if err := o.ParseFrom(osName, strings.EqualFold); err != nil {
				t.Fatalf("OS.ParseFrom(%q) unexpected error: %v", osName, err)
			}

			var e platform.Extension

			e.ParseFrom(o)

			if !e.IsNil() {
				t.Errorf("%s Extension.IsNil() = false, want true (got %q)", osName, e.String())
			}
		})
	}
}

func TestExtensionNonNil(t *testing.T) {
	t.Parallel()

	t.Run("directly assigned non-empty extension is not nil", func(t *testing.T) {
		t.Parallel()

		e := platform.Extension(".tar.gz")

		if e.IsNil() {
			t.Error("non-empty Extension.IsNil() = true, want false")
		}

		if e.String() != ".tar.gz" {
			t.Errorf("String() = %q, want %q", e.String(), ".tar.gz")
		}
	})

	t.Run("exe extension is not nil", func(t *testing.T) {
		t.Parallel()

		e := platform.Extension(".exe")

		if e.IsNil() {
			t.Error(".exe Extension.IsNil() = true, want false")
		}
	})
}
