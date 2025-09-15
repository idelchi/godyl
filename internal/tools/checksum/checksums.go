package checksum

import (
	"strings"

	"github.com/idelchi/godyl/internal/debug"
	"github.com/idelchi/godyl/pkg/score"
)

// Checksums represents a list of checksum file names or URLs.
type Checksums []string

// Indicators returns a list of common substrings found in checksum file names.
func Indicators() []string {
	return []string{
		"checksum",
		"checksums",
		"sha256",
		"sha-256",
		"sha512",
		"sha-512",
		"md5",
		"md-5",
		"shasums",
		"sha1",
		"sha-1",
		"hash",
		"hashes",
		"digest",
		"digests",
		"sum",
		"sums",
		"shasum",
		"sha256sum",
		"sha256sums",
		"sha512sum",
		"sha512sums",
		"sha1sum",
		"sha1sums",
		"md5sum",
		"md5sums",
		"sha",
	}
}

// IsChecksumLike returns the checksums that appear to be checksum files based on common indicators.
func (cs Checksums) IsChecksumLike() Checksums {
	var found Checksums

	for _, name := range cs {
		for _, indicator := range Indicators() {
			debug.Debug("checking if %q contains %q", strings.ToLower(name), indicator)

			if strings.Contains(strings.ToLower(name), indicator) {
				debug.Debug("found checksum-like file: %q", name)

				found = append(found, name)

				break
			}
		}
	}

	return found
}

// Preferred returns the likeliest checksum given a name.
func (cs Checksums) Preferred(name string) string {
	if len(cs) == 1 {
		return cs[0]
	}

	if len(cs) == 0 {
		return ""
	}

	lowerName := strings.ToLower(name)

	conditions := []func(string) int{
		func(c string) int {
			if strings.Contains(strings.ToLower(c), lowerName) {
				return 2 //nolint:mnd	// Magic number is acceptable here
			}

			return 0
		},
	}

	patterns := map[string]int{
		"checksum":      1,
		"checksum.txt":  1,
		"checksums.txt": 1,
	}

	for p, w := range patterns {
		cond := func(p string, w int) func(string) int {
			return func(c string) int {
				if strings.Contains(strings.ToLower(c), p) {
					return w
				}

				return 0
			}
		}

		conditions = append(conditions, cond(p, w))
	}

	scores := score.Score(cs, conditions...)

	top := scores.Top()

	debug.Debug("all checksum candidates for %q: %v", name, scores)
	debug.Debug("best checksum candidates for %q: %v", name, top)

	best := top[0]

	if best.Score == 0 {
		return ""
	}

	return best.Item
}
