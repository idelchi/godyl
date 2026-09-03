package processor

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"sync"

	"golang.org/x/sync/singleflight"

	"github.com/idelchi/godyl/internal/data"
	"github.com/idelchi/godyl/internal/tools/mode"
	"github.com/idelchi/godyl/internal/tools/result"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/internal/tools/sources/install"
	"github.com/idelchi/godyl/internal/tools/tool"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/folder"
)

// acquiredArtifact tracks a downloaded artifact and the temp dir that owns it.
type acquiredArtifact struct {
	// destination is the downloaded artifact shared by matching tools.
	destination file.File
	// dir is the temporary directory that must be removed after processing.
	dir folder.Folder
}

// artifactStore reuses resolved artifacts within a single processing run.
type artifactStore struct {
	// group coalesces concurrent downloads for the same artifact key.
	group singleflight.Group
	// artifacts caches completed downloads by artifact key.
	artifacts map[string]acquiredArtifact
	// mu protects artifacts.
	mu sync.Mutex
}

// artifactKeyData contains the stable inputs that affect downloaded bytes.
type artifactKeyData struct {
	// Headers contains normalized request headers.
	Headers map[string][]string `json:"headers"`
	// ChecksumQuery contains the source-specific checksum lookup expression.
	ChecksumQuery string `json:"checksum_query"`
	// ChecksumType names the digest algorithm.
	ChecksumType string `json:"checksum_type"`
	// ChecksumValue is the expected digest.
	ChecksumValue string `json:"checksum_value"`
	// NoVerifyChecksum records whether checksum validation is disabled.
	NoVerifyChecksum bool `json:"no_verify_checksum"`
	// NoVerifySSL records whether TLS certificate validation is disabled.
	NoVerifySSL bool `json:"no_verify_ssl"`
	// Source identifies the release provider.
	Source string `json:"source"`
	// URL is the resolved artifact URL.
	URL string `json:"url"`
}

// newArtifactStore creates an empty per-run artifact store.
func newArtifactStore() *artifactStore {
	return &artifactStore{
		artifacts: make(map[string]acquiredArtifact),
	}
}

// acquire downloads an artifact once and returns the shared destination.
func (s *artifactStore) acquire(source sources.Type, d install.Data) (file.File, error) {
	key, err := newArtifactKey(source, d)
	if err != nil {
		return "", err
	}

	s.mu.Lock()

	if artifact, ok := s.artifacts[key]; ok {
		s.mu.Unlock()

		return artifact.destination, nil
	}

	s.mu.Unlock()

	value, err, _ := s.group.Do(key, func() (any, error) {
		s.mu.Lock()

		if artifact, ok := s.artifacts[key]; ok {
			s.mu.Unlock()

			return artifact.destination, nil
		}

		s.mu.Unlock()

		dir, err := data.CreateUniqueDirIn()
		if err != nil {
			return file.File(""), fmt.Errorf("creating random dir: %w", err)
		}

		destination, err := install.Acquire(d, dir)
		if err != nil {
			return file.File(""), errors.Join(err, dir.Remove())
		}

		s.mu.Lock()

		s.artifacts[key] = acquiredArtifact{
			destination: destination,
			dir:         dir,
		}

		s.mu.Unlock()

		return destination, nil
	})
	if err != nil {
		return "", err
	}

	destination, ok := value.(file.File)
	if !ok {
		return "", fmt.Errorf("artifact store returned %T", value)
	}

	return destination, nil
}

// cleanup removes temp dirs owned by shared artifacts.
func (s *artifactStore) cleanup() error {
	s.mu.Lock()

	artifacts := make([]acquiredArtifact, 0, len(s.artifacts))
	for _, artifact := range s.artifacts {
		artifacts = append(artifacts, artifact)
	}

	s.mu.Unlock()

	var err error

	for _, artifact := range artifacts {
		err = errors.Join(err, artifact.dir.Remove())
	}

	return err
}

// downloadTool installs a tool, reusing its artifact when it is safe to share.
func (p *Processor) downloadTool(ctx context.Context, t *tool.Tool, artifacts *artifactStore) result.Result {
	if !isReusableArtifact(t) {
		return t.Download(ctx, p.progress.Tracker())
	}

	data, err := t.InstallData()
	if err != nil {
		return result.WithFailed(err.Error())
	}

	data.ProgressListener = p.progress.Tracker()

	destination, err := artifacts.acquire(t.Source.Type, data)
	if err != nil {
		return result.WithFailed("installing tool").Wrap(err)
	}

	if _, err := install.Find(destination, data); err != nil {
		return result.WithFailed("installing tool").Wrap(err)
	}

	return t.CompleteInstall(ctx)
}

// isReusableArtifact reports whether a tool can share its resolved artifact.
func isReusableArtifact(t *tool.Tool) bool {
	return t.Source.Type.SupportsArtifactReuse() &&
		t.Mode == mode.Find &&
		t.URL != ""
}

// newArtifactKey hashes the download inputs that identify one reusable artifact.
func newArtifactKey(source sources.Type, d install.Data) (string, error) {
	keyData := artifactKeyData{
		Headers:          normalizeHeaders(d.Header),
		ChecksumQuery:    d.Checksum.ToQuery(),
		ChecksumType:     d.Checksum.Type.String(),
		ChecksumValue:    d.Checksum.Value,
		NoVerifyChecksum: d.NoVerifyChecksum,
		NoVerifySSL:      d.NoVerifySSL,
		Source:           source.String(),
		URL:              d.Path,
	}

	material, err := json.Marshal(keyData)
	if err != nil {
		return "", fmt.Errorf("marshaling artifact key: %w", err)
	}

	sum := sha256.Sum256(material)

	return hex.EncodeToString(sum[:]), nil
}

// normalizeHeaders canonicalizes headers for stable artifact keys.
func normalizeHeaders(headers http.Header) map[string][]string {
	normalized := make(map[string][]string)

	for key, values := range headers {
		canonical := http.CanonicalHeaderKey(key)

		normalized[canonical] = append(normalized[canonical], values...)
	}

	for key := range normalized {
		sort.Strings(normalized[key])
	}

	return normalized
}
