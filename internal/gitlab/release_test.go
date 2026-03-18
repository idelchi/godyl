package gitlab_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"

	internalgitlab "github.com/idelchi/godyl/internal/gitlab"
	"github.com/idelchi/godyl/internal/release"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func TestFromRepositoryRelease(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     *gitlab.Release
		want      *release.Release
		wantErrIs error
	}{
		{
			name:      "nil input returns ErrRelease",
			input:     nil,
			wantErrIs: release.ErrRelease,
		},
		{
			name: "empty tag name returns ErrRelease",
			input: &gitlab.Release{
				Name:    "Some Release",
				TagName: "",
			},
			wantErrIs: release.ErrRelease,
		},
		{
			name: "basic field mapping",
			input: &gitlab.Release{
				Name:    "v1.2.3 Release",
				TagName: "v1.2.3",
				Assets:  gitlab.ReleaseAssets{Links: []*gitlab.ReleaseLink{}},
			},
			want: &release.Release{
				Name:   "v1.2.3 Release",
				Tag:    "v1.2.3",
				Assets: release.Assets{},
			},
		},
		{
			name: "nil links produces empty assets",
			input: &gitlab.Release{
				Name:    "empty",
				TagName: "v0.1.0",
				Assets:  gitlab.ReleaseAssets{Links: nil},
			},
			want: &release.Release{
				Name:   "empty",
				Tag:    "v0.1.0",
				Assets: release.Assets{},
			},
		},
		{
			name: "single package link",
			input: &gitlab.Release{
				Name:    "Test Release",
				TagName: "v1.0.0",
				Assets: gitlab.ReleaseAssets{
					Links: []*gitlab.ReleaseLink{
						{
							Name:           "godyl_linux_amd64.tar.gz",
							URL:            "https://example.com/files/godyl_linux_amd64.tar.gz",
							DirectAssetURL: "https://example.com/direct/godyl_linux_amd64.tar.gz",
							LinkType:       gitlab.PackageLinkType,
						},
					},
				},
			},
			want: &release.Release{
				Name: "Test Release",
				Tag:  "v1.0.0",
				Assets: release.Assets{
					{
						Name: "godyl_linux_amd64.tar.gz",
						URL:  "https://example.com/direct/godyl_linux_amd64.tar.gz",
						Type: "package",
					},
				},
			},
		},
		{
			name: "multiple links with different types",
			input: &gitlab.Release{
				Name:    "Test Release",
				TagName: "v1.0.0",
				Assets: gitlab.ReleaseAssets{
					Links: []*gitlab.ReleaseLink{
						{
							Name:           "godyl_windows_amd64.zip",
							URL:            "https://example.com/files/godyl_windows_amd64.zip",
							DirectAssetURL: "https://example.com/direct/godyl_windows_amd64.zip",
							LinkType:       gitlab.PackageLinkType,
						},
						{
							Name:           "checksums.txt",
							URL:            "https://example.com/files/checksums.txt",
							DirectAssetURL: "https://example.com/direct/checksums.txt",
							LinkType:       gitlab.OtherLinkType,
						},
						{
							Name:           "godyl_docker.tar",
							URL:            "https://example.com/files/godyl_docker.tar",
							DirectAssetURL: "https://example.com/direct/godyl_docker.tar",
							LinkType:       gitlab.ImageLinkType,
						},
					},
				},
			},
			want: &release.Release{
				Name: "Test Release",
				Tag:  "v1.0.0",
				Assets: release.Assets{
					{
						Name: "godyl_windows_amd64.zip",
						URL:  "https://example.com/direct/godyl_windows_amd64.zip",
						Type: "package",
					},
					{Name: "checksums.txt", URL: "https://example.com/direct/checksums.txt", Type: "other"},
					{Name: "godyl_docker.tar", URL: "https://example.com/direct/godyl_docker.tar", Type: "image"},
				},
			},
		},
		{
			name: "runbook link type",
			input: &gitlab.Release{
				Name:    "Runbook Release",
				TagName: "v3.0.0",
				Assets: gitlab.ReleaseAssets{
					Links: []*gitlab.ReleaseLink{
						{
							Name:           "runbook.pdf",
							URL:            "https://example.com/runbook.pdf",
							DirectAssetURL: "https://example.com/direct/runbook.pdf",
							LinkType:       gitlab.RunbookLinkType,
						},
					},
				},
			},
			want: &release.Release{
				Name: "Runbook Release",
				Tag:  "v3.0.0",
				Assets: release.Assets{
					{Name: "runbook.pdf", URL: "https://example.com/direct/runbook.pdf", Type: "runbook"},
				},
			},
		},
		{
			name: "description field is not mapped to Body",
			input: &gitlab.Release{
				Name:        "Described Release",
				TagName:     "v4.0.0",
				Description: "This is the release description",
				Assets:      gitlab.ReleaseAssets{Links: []*gitlab.ReleaseLink{}},
			},
			want: &release.Release{
				Name:   "Described Release",
				Tag:    "v4.0.0",
				Body:   "",
				Assets: release.Assets{},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := internalgitlab.FromRepositoryRelease(tc.input)

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

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("FromRepositoryRelease() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestFromRepositoryRelease_NilAssetLink(t *testing.T) {
	t.Parallel()

	// FromRepositoryRelease iterates repoRelease.Assets.Links without nil-checking
	// each element. A nil *ReleaseLink in the slice causes a nil-pointer dereference
	// when the loop body accesses link.Name, link.DirectAssetURL, or link.LinkType.
	// This test documents that current (unfixed) behaviour.
	input := &gitlab.Release{
		Name:    "Test Release",
		TagName: "v1.0.0",
		Assets: gitlab.ReleaseAssets{
			Links: []*gitlab.ReleaseLink{
				nil, // nil entry — triggers a panic in the current implementation
				{
					Name:           "tool-linux-amd64.tar.gz",
					URL:            "https://example.com/files/tool-linux-amd64.tar.gz",
					DirectAssetURL: "https://example.com/direct/tool-linux-amd64.tar.gz",
					LinkType:       gitlab.PackageLinkType,
				},
			},
		},
	}

	defer func() {
		r := recover()
		if r == nil {
			t.Error("expected a panic from nil *ReleaseLink dereference, but no panic occurred")
		}
	}()

	_, _ = internalgitlab.FromRepositoryRelease(input)
}

func TestFromRepositoryRelease_EmptyDirectAssetURL(t *testing.T) {
	t.Parallel()

	// A ReleaseLink with an empty DirectAssetURL is valid from the struct perspective.
	// FromRepositoryRelease maps DirectAssetURL directly to Asset.URL, so an empty
	// URL propagates through without error.
	input := &gitlab.Release{
		Name:    "Release with empty URL",
		TagName: "v1.0.0",
		Assets: gitlab.ReleaseAssets{
			Links: []*gitlab.ReleaseLink{
				{
					Name:           "tool-linux-amd64.tar.gz",
					URL:            "https://example.com/files/tool-linux-amd64.tar.gz",
					DirectAssetURL: "", // empty direct asset URL
					LinkType:       gitlab.PackageLinkType,
				},
			},
		},
	}

	got, err := internalgitlab.FromRepositoryRelease(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := &release.Release{
		Name: "Release with empty URL",
		Tag:  "v1.0.0",
		Assets: release.Assets{
			{
				Name: "tool-linux-amd64.tar.gz",
				URL:  "", // empty DirectAssetURL propagates as-is
				Type: "package",
			},
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("FromRepositoryRelease_EmptyDirectAssetURL() mismatch (-want +got):\n%s", diff)
	}
}
