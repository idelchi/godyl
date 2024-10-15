package goi

import (
	"fmt"

	"github.com/idelchi/godyl/internal/detect"
	"github.com/idelchi/godyl/internal/match"
)

type Targets struct {
	Files []Target `json:"files"`
}

func (gt Targets) FilterBy(predicate func(Target) bool) Targets {
	var filtered Targets
	for _, file := range gt.Files {
		if predicate(file) {
			filtered.Files = append(filtered.Files, file)
		}
	}
	return filtered
}

func (gt Targets) FilterByOS(os string) Targets {
	return gt.FilterBy(func(file Target) bool {
		return file.OS == os
	})
}

func (gt Targets) FilterByArch(arch string) Targets {
	return gt.FilterBy(func(file Target) bool {
		return file.Arch == arch
	})
}

func (t Targets) Match() (match.Results, error) {
	platform := detect.Platform{}
	if err := platform.Detect(); err != nil {
		return nil, fmt.Errorf("detecting platform: %w", err)
	}

	var assets match.Assets

	for _, tt := range t.Files {
		asset := match.Asset{Name: tt.FileName}

		asset.Platform.OS.From(tt.OS)
		asset.Platform.Architecture.From(tt.Arch, "")

		assets = append(assets, asset)
	}

	hints := []match.Hint{
		match.NewDefaultHint(platform.OS.Name()),
	}

	var err error

	matches := assets.Select(match.Requirements{Platform: platform, Hints: hints})
	switch {
	case !matches.HasQualified():
		err = fmt.Errorf("no qualified file found")
	case matches.IsAmbigious():
		err = fmt.Errorf("ambiguous file selection")
	case !matches.Success():
		err = fmt.Errorf("no matching file found")
	}

	return matches, err
}
