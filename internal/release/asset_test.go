package release_test

import (
	"slices"
	"testing"

	"github.com/idelchi/godyl/internal/release"
)

func TestAssetMatch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		assetName string
		pattern   string
		wantMatch bool
		wantErr   bool
	}{
		{
			name:      "glob pattern matches tar.gz",
			assetName: "tool.tar.gz",
			pattern:   "*.tar.gz",
			wantMatch: true,
		},
		{
			name:      "zip pattern does not match tar.gz",
			assetName: "tool.tar.gz",
			pattern:   "*.zip",
			wantMatch: false,
		},
		{
			name:      "exact name matches itself",
			assetName: "mytool-v1.0.0-linux-amd64",
			pattern:   "mytool-v1.0.0-linux-amd64",
			wantMatch: true,
		},
		{
			name:      "exact name does not match different name",
			assetName: "mytool-v1.0.0-linux-amd64",
			pattern:   "mytool-v2.0.0-linux-amd64",
			wantMatch: false,
		},
		{
			name:      "invalid pattern returns error",
			assetName: "tool.tar.gz",
			pattern:   "[invalid",
			wantErr:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			a := release.Asset{Name: tc.assetName}
			got, err := a.Match(tc.pattern)

			if tc.wantErr {
				if err == nil {
					t.Errorf("Match(%q) expected error, got nil", tc.pattern)
				}

				return
			}

			if err != nil {
				t.Fatalf("Match(%q) unexpected error: %v", tc.pattern, err)
			}

			if got != tc.wantMatch {
				t.Errorf("Match(%q) = %v, want %v", tc.pattern, got, tc.wantMatch)
			}
		})
	}
}

func TestAssetHasExtension(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		assetName string
		extension string
		wantMatch bool
	}{
		{
			name:      "single dot extension matches",
			assetName: "file.tar.gz",
			extension: ".gz",
			wantMatch: true,
		},
		{
			name:      "compound extension matched as full suffix",
			assetName: "file.tar.gz",
			extension: ".tar.gz",
			wantMatch: true,
		},
		{
			name:      "different single dot extension does not match",
			assetName: "file.tar.gz",
			extension: ".zip",
			wantMatch: false,
		},
		{
			name:      "exe extension matches exe file",
			assetName: "tool.exe",
			extension: ".exe",
			wantMatch: true,
		},
		{
			name:      "multi-dot extension does not match wrong suffix",
			assetName: "file.tar.gz",
			extension: ".tar.bz2",
			wantMatch: false,
		},
		{
			name:      "no extension does not match extension",
			assetName: "toolbinary",
			extension: ".gz",
			wantMatch: false,
		},
		{
			// filepath.Ext("toolbinary") == "" and extension == "", so "" == "" is true.
			name:      "empty extension matches file with no extension",
			assetName: "toolbinary",
			extension: "",
			wantMatch: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			a := release.Asset{Name: tc.assetName}

			got, err := a.HasExtension(tc.extension)
			if err != nil {
				t.Fatalf("HasExtension(%q) unexpected error: %v", tc.extension, err)
			}

			if got != tc.wantMatch {
				t.Errorf("HasExtension(%q) = %v, want %v", tc.extension, got, tc.wantMatch)
			}
		})
	}
}

func TestAssetsFilterByName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		assets    release.Assets
		filterBy  string
		wantNames []string
	}{
		{
			name: "case-insensitive match finds both variants",
			assets: release.Assets{
				{Name: "tool-linux-amd64.tar.gz"},
				{Name: "tool-darwin-amd64.tar.gz"},
				{Name: "Tool-linux-amd64.tar.gz"},
				{Name: "other-linux-amd64.tar.gz"},
			},
			filterBy:  "TOOL-LINUX-AMD64.TAR.GZ",
			wantNames: []string{"Tool-linux-amd64.tar.gz", "tool-linux-amd64.tar.gz"},
		},
		{
			name: "darwin variant only",
			assets: release.Assets{
				{Name: "tool-linux-amd64.tar.gz"},
				{Name: "tool-darwin-amd64.tar.gz"},
				{Name: "Tool-linux-amd64.tar.gz"},
				{Name: "other-linux-amd64.tar.gz"},
			},
			filterBy:  "tool-darwin-amd64.tar.gz",
			wantNames: []string{"tool-darwin-amd64.tar.gz"},
		},
		{
			name: "no match returns empty",
			assets: release.Assets{
				{Name: "tool-linux-amd64.tar.gz"},
				{Name: "tool-darwin-amd64.tar.gz"},
				{Name: "Tool-linux-amd64.tar.gz"},
				{Name: "other-linux-amd64.tar.gz"},
			},
			filterBy:  "nonexistent.tar.gz",
			wantNames: nil,
		},
		{
			name:      "nil assets returns nil",
			assets:    nil,
			filterBy:  "tool-linux-amd64.tar.gz",
			wantNames: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := tc.assets.FilterByName(tc.filterBy)

			gotNames := make([]string, len(got))
			for i, a := range got {
				gotNames[i] = a.Name
			}

			slices.Sort(gotNames)

			wantSorted := slices.Clone(tc.wantNames)
			slices.Sort(wantSorted)

			if !slices.Equal(gotNames, wantSorted) {
				t.Errorf("FilterByName(%q) = %v, want %v", tc.filterBy, gotNames, wantSorted)
			}
		})
	}
}

func TestAssetsChecksums(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		assets  release.Assets
		pattern string
		want    []string
	}{
		{
			name:    "empty assets returns nil",
			assets:  nil,
			pattern: "",
			want:    nil,
		},
		{
			name: "empty pattern returns all checksum-like assets",
			assets: release.Assets{
				{Name: "tool-linux-amd64.tar.gz"},
				{Name: "tool-darwin-amd64.tar.gz"},
				{Name: "checksums.txt"},
				{Name: "tool_sha256sums.txt"},
				{Name: "tool-linux-amd64.tar.gz.md5"},
			},
			pattern: "",
			want:    []string{"checksums.txt", "tool-linux-amd64.tar.gz.md5", "tool_sha256sums.txt"},
		},
		{
			name: "pattern filters checksum assets by glob",
			assets: release.Assets{
				{Name: "tool-linux-amd64.tar.gz"},
				{Name: "tool-darwin-amd64.tar.gz"},
				{Name: "checksums.txt"},
				{Name: "tool_sha256sums.txt"},
				{Name: "tool-linux-amd64.tar.gz.md5"},
			},
			pattern: "checksums*",
			want:    []string{"checksums.txt"},
		},
		{
			name: "pattern with no match returns nil",
			assets: release.Assets{
				{Name: "tool-linux-amd64.tar.gz"},
				{Name: "checksums.txt"},
			},
			pattern: "*.sig",
			want:    nil,
		},
		{
			name: "pattern matches multiple checksum files",
			assets: release.Assets{
				{Name: "tool-linux-amd64.tar.gz"},
				{Name: "tool-darwin-amd64.tar.gz"},
				{Name: "checksums.txt"},
				{Name: "tool_sha256sums.txt"},
				{Name: "tool-linux-amd64.tar.gz.md5"},
			},
			pattern: "tool*",
			want:    []string{"tool-linux-amd64.tar.gz.md5", "tool_sha256sums.txt"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := tc.assets.Checksums(tc.pattern)

			gotSorted := slices.Clone([]string(got))
			slices.Sort(gotSorted)

			wantSorted := slices.Clone(tc.want)
			slices.Sort(wantSorted)

			if !slices.Equal(gotSorted, wantSorted) {
				t.Errorf("Checksums(%q) = %v, want %v", tc.pattern, gotSorted, wantSorted)
			}
		})
	}
}
