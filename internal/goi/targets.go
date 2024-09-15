package goi

import (
	"errors"
	"fmt"

	"github.com/idelchi/godyl/internal/detect"
	"github.com/idelchi/godyl/internal/match"
)

// Targets represents a collection of Target files associated with a Go release.
type Targets struct {
	Files []Target `json:"files"` // Files is the list of Target files available in the release.
}

// FilterBy filters the Targets based on a given predicate function. It returns a new Targets collection containing
// only the files that match the provided condition.
func (gt Targets) FilterBy(predicate func(Target) bool) Targets {
	var filtered Targets

	for _, file := range gt.Files {
		if predicate(file) {
			filtered.Files = append(filtered.Files, file)
		}
	}

	return filtered
}

// FilterByOS filters the Targets to include only those files that match the specified operating system (OS).
func (gt Targets) FilterByOS(os string) Targets {
	return gt.FilterBy(func(file Target) bool {
		return file.OS == os
	})
}

// FilterByArch filters the Targets to include only those files that match the specified architecture.
func (gt Targets) FilterByArch(arch string) Targets {
	return gt.FilterBy(func(file Target) bool {
		return file.Arch == arch
	})
}

// Match attempts to find the best matching file from the Targets collection based on the platform detected
// by the system. It returns a list of matched results or an error if no suitable match is found.
func (gt Targets) Match() (match.Results, error) {
	platform := detect.Platform{}
	if err := platform.Detect(); err != nil {
		return nil, fmt.Errorf("detecting platform: %w", err)
	}

	var assets match.Assets

	for _, tt := range gt.Files {
		asset := match.Asset{Name: tt.FileName}

		asset.Platform.OS.ParseFrom(tt.OS)             //nolint:errcheck
		asset.Platform.Architecture.ParseFrom(tt.Arch) //nolint:errcheck

		assets = append(assets, asset)
	}

	hints := []match.Hint{
		{
			Pattern: platform.OS.String(),
			Must:    true,
		},
	}

	var err error

	matches := assets.Select(match.Requirements{Platform: platform, Hints: hints})

	switch {
	case !matches.HasQualified():
		err = fmt.Errorf("%w: no qualified file found", ErrMatch)
	case matches.IsAmbigious():
		err = fmt.Errorf("%w: ambiguous file selection", ErrMatch)
	case !matches.Success():
		err = fmt.Errorf("%w: no matching file found", ErrMatch)
	}

	return matches, err
}

// ErrMatch is returned when no suitable match is found.
var ErrMatch = errors.New("matching")
