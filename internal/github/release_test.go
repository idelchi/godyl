package github_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	gogithub "github.com/google/go-github/v74/github"

	internalgithub "github.com/idelchi/godyl/internal/github"
	"github.com/idelchi/godyl/internal/release"
)

func TestFromRepositoryRelease(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     *gogithub.RepositoryRelease
		want      *release.Release
		wantErrIs error
	}{
		{
			name:      "nil release",
			input:     nil,
			wantErrIs: release.ErrRelease,
		},
		{
			name: "nil tag name",
			input: &gogithub.RepositoryRelease{
				Name: gogithub.Ptr("Release with no tag"),
				Body: gogithub.Ptr("Some body"),
			},
			wantErrIs: release.ErrRelease,
		},
		{
			name: "full release with two assets",
			input: &gogithub.RepositoryRelease{
				Name:    gogithub.Ptr("Release v1.2.3"),
				TagName: gogithub.Ptr("v1.2.3"),
				Body:    gogithub.Ptr("Release notes body"),
				Assets: []*gogithub.ReleaseAsset{
					{
						Name: gogithub.Ptr("tool-linux-amd64.tar.gz"),
						BrowserDownloadURL: gogithub.Ptr(
							"https://github.com/owner/repo/releases/download/v1.2.3/tool-linux-amd64.tar.gz",
						),
						ContentType: gogithub.Ptr("application/gzip"),
						Digest:      gogithub.Ptr("sha256:abc123def456"),
					},
					{
						Name: gogithub.Ptr("tool-darwin-arm64.tar.gz"),
						BrowserDownloadURL: gogithub.Ptr(
							"https://github.com/owner/repo/releases/download/v1.2.3/tool-darwin-arm64.tar.gz",
						),
						ContentType: gogithub.Ptr("application/gzip"),
						Digest:      gogithub.Ptr("sha256:deadbeef1234"),
					},
				},
			},
			want: &release.Release{
				Name: "Release v1.2.3",
				Tag:  "v1.2.3",
				Body: "Release notes body",
				Assets: release.Assets{
					{
						Name:   "tool-linux-amd64.tar.gz",
						URL:    "https://github.com/owner/repo/releases/download/v1.2.3/tool-linux-amd64.tar.gz",
						Type:   "application/gzip",
						Digest: "sha256:abc123def456",
					},
					{
						Name:   "tool-darwin-arm64.tar.gz",
						URL:    "https://github.com/owner/repo/releases/download/v1.2.3/tool-darwin-arm64.tar.gz",
						Type:   "application/gzip",
						Digest: "sha256:deadbeef1234",
					},
				},
			},
		},
		{
			name: "empty assets list",
			input: &gogithub.RepositoryRelease{
				Name:    gogithub.Ptr("Release v2.0.0"),
				TagName: gogithub.Ptr("v2.0.0"),
				Body:    gogithub.Ptr(""),
				Assets:  []*gogithub.ReleaseAsset{},
			},
			want: &release.Release{
				Tag:    "v2.0.0",
				Name:   "Release v2.0.0",
				Assets: release.Assets{},
			},
		},
		{
			name: "nil optional fields",
			input: &gogithub.RepositoryRelease{
				TagName: gogithub.Ptr("v1.0.0"),
			},
			want: &release.Release{
				Tag:    "v1.0.0",
				Assets: release.Assets{},
			},
		},
		{
			name: "asset with nil Name is skipped",
			input: &gogithub.RepositoryRelease{
				TagName: gogithub.Ptr("v1.0.0"),
				Assets: []*gogithub.ReleaseAsset{
					{
						Name:               nil,
						BrowserDownloadURL: gogithub.Ptr("https://example.com/file.tar.gz"),
						ContentType:        gogithub.Ptr("application/gzip"),
					},
				},
			},
			want: &release.Release{
				Tag:    "v1.0.0",
				Assets: release.Assets{},
			},
		},
		{
			name: "asset with nil BrowserDownloadURL is skipped",
			input: &gogithub.RepositoryRelease{
				TagName: gogithub.Ptr("v1.0.0"),
				Assets: []*gogithub.ReleaseAsset{
					{
						Name:               gogithub.Ptr("file.tar.gz"),
						BrowserDownloadURL: nil,
						ContentType:        gogithub.Ptr("application/gzip"),
					},
				},
			},
			want: &release.Release{
				Tag:    "v1.0.0",
				Assets: release.Assets{},
			},
		},
		{
			name: "asset with nil ContentType is skipped",
			input: &gogithub.RepositoryRelease{
				TagName: gogithub.Ptr("v1.0.0"),
				Assets: []*gogithub.ReleaseAsset{
					{
						Name:               gogithub.Ptr("file.tar.gz"),
						BrowserDownloadURL: gogithub.Ptr("https://example.com/file.tar.gz"),
						ContentType:        nil,
					},
				},
			},
			want: &release.Release{
				Tag:    "v1.0.0",
				Assets: release.Assets{},
			},
		},
		{
			name: "nil asset pointer in Assets slice does not panic",
			input: &gogithub.RepositoryRelease{
				TagName: gogithub.Ptr("v1.0.0"),
				Assets: []*gogithub.ReleaseAsset{
					nil,
					{
						Name: gogithub.Ptr("tool.tar.gz"),
						BrowserDownloadURL: gogithub.Ptr(
							"https://github.com/owner/repo/releases/download/v1.0.0/tool.tar.gz",
						),
						ContentType: gogithub.Ptr("application/gzip"),
					},
				},
			},
			want: &release.Release{
				Tag: "v1.0.0",
				Assets: release.Assets{
					{
						Name: "tool.tar.gz",
						URL:  "https://github.com/owner/repo/releases/download/v1.0.0/tool.tar.gz",
						Type: "application/gzip",
					},
				},
			},
		},
		{
			name: "nil asset digest produces empty Digest field",
			input: &gogithub.RepositoryRelease{
				TagName: gogithub.Ptr("v1.0.0"),
				Assets: []*gogithub.ReleaseAsset{
					{
						Name: gogithub.Ptr("tool.tar.gz"),
						BrowserDownloadURL: gogithub.Ptr(
							"https://github.com/owner/repo/releases/download/v1.0.0/tool.tar.gz",
						),
						ContentType: gogithub.Ptr("application/gzip"),
						Digest:      nil,
					},
				},
			},
			want: &release.Release{
				Tag: "v1.0.0",
				Assets: release.Assets{
					{
						Name: "tool.tar.gz",
						URL:  "https://github.com/owner/repo/releases/download/v1.0.0/tool.tar.gz",
						Type: "application/gzip",
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := internalgithub.FromRepositoryRelease(tc.input)

			if tc.wantErrIs != nil {
				if err == nil {
					t.Fatalf("expected error wrapping %v, got nil", tc.wantErrIs)
				}

				if !errors.Is(err, tc.wantErrIs) {
					t.Errorf("expected errors.Is(err, %v), got %v", tc.wantErrIs, err)
				}

				if got != nil {
					t.Errorf("expected nil result on error, got %+v", got)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got == nil {
				t.Fatal("expected non-nil release")
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("FromRepositoryRelease() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
