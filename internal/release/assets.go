package release

import (
	"context"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/idelchi/godyl/internal/detect/platform"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/checksum"
)

// Assets represents a collection of release assets.
type Assets []Asset

// Select returns the single deterministic match, or delegates an ambiguous
// set of equally ranked candidates to the optional selector.
func (as Assets) Select(ctx context.Context, requirements match.Requirements) (Asset, error) {
	matches := as.Match(requirements)

	if matches.HasErrors() {
		return Asset{}, matches.Errors()[0]
	}

	matchErr := matches.Status()
	if matchErr != nil {
		matchErr = matches.WithoutZero().Status()
	}

	if matchErr == nil {
		return as.assetByName(matches[0].Asset.Name)
	}

	diagnosticRequirements := requirements

	diagnosticRequirements.Selector = nil

	selectionErr := &match.SelectionError{
		Cause:        matchErr,
		Candidates:   as.names(),
		Results:      matches,
		Requirements: diagnosticRequirements,
	}

	if requirements.Selector == nil || !errors.Is(matchErr, match.ErrAmbiguous) {
		return Asset{}, selectionErr
	}

	candidates := make([]string, len(matches))
	for i, result := range matches {
		candidates[i] = result.Asset.Name
	}

	selected, err := requirements.Selector.Select(ctx, match.SelectionRequest{
		Target:     requirements.Target,
		Platform:   requirements.Platform,
		Candidates: candidates,
	})
	if err != nil {
		return Asset{}, errors.Join(selectionErr, fmt.Errorf("AI fallback: %w", err))
	}

	if !slices.Contains(candidates, selected) {
		return Asset{}, errors.Join(
			selectionErr,
			fmt.Errorf("AI fallback: selected asset %q is not an ambiguous candidate", selected),
		)
	}

	asset, err := as.assetByName(selected)
	if err != nil {
		return Asset{}, errors.Join(selectionErr, fmt.Errorf("AI fallback: %w", err))
	}

	return asset, nil
}

// FilterByName returns the assets that match the given name.
// It compares asset names in a case-insensitive manner.
func (as Assets) FilterByName(name string) (assets Assets) {
	for _, asset := range as {
		if strings.EqualFold(asset.Name, name) {
			assets = append(assets, asset)
		}
	}

	return assets
}

// Match checks if the assets match the given requirements.
// It processes each asset to extract platform and extension information.
func (as Assets) Match(requirements match.Requirements) (matches match.Results) {
	assets := make(match.Assets, 0, len(as))

	for _, a := range as {
		asset := match.Asset{Name: a.Name}
		asset.Parse()

		asset.Platform.Extension = platform.Extension(
			filepath.Ext(a.Name),
		)

		assets = append(assets, asset)
	}

	matches = assets.Select(requirements)

	return matches
}

// Checksums returns all assets that appear to be checksum files.
// When pattern is non-empty, only checksum-like assets matching the pattern are returned.
func (as Assets) Checksums(pattern string) checksum.Checksums {
	names := make(checksum.Checksums, 0, len(as))

	for _, asset := range as {
		names = append(names, asset.Name)
	}

	checksums := names.IsChecksumLike()

	if pattern == "" {
		return checksums
	}

	var filtered checksum.Checksums

	for _, cs := range checksums {
		matched, err := path.Match(pattern, cs)
		if err == nil && matched {
			filtered = append(filtered, cs)
		}
	}

	return filtered
}

func (as Assets) assetByName(name string) (Asset, error) {
	assets := as.FilterByName(name)
	if len(assets) != 1 {
		return Asset{}, fmt.Errorf("expected one release asset named %q, found %d", name, len(assets))
	}

	return assets[0], nil
}

func (as Assets) names() []string {
	names := make([]string, len(as))
	for i, asset := range as {
		names[i] = asset.Name
	}

	return names
}
