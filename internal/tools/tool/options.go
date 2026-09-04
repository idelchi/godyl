package tool

import "github.com/idelchi/godyl/internal/match"

// ResolveOption is a functional option type for the Resolve method.
type ResolveOption func(*resolveOptions)

// resolveOptions holds all configurable options for the Resolve method.
type resolveOptions struct {
	// assetSelector is consulted only for ambiguous deterministic asset matches.
	assetSelector match.AssetSelector
	// skipVersion bypasses version resolution.
	skipVersion bool
	// upUntilVersion stops resolution after determining the version.
	upUntilVersion bool
	// skipURL bypasses artifact URL resolution.
	skipURL bool
}

// WithAssetSelector enables a tie-breaker for ambiguous release assets.
func WithAssetSelector(selector match.AssetSelector) ResolveOption {
	return func(o *resolveOptions) {
		o.assetSelector = selector
	}
}

// WithoutVersion returns a ResolveOption that skips version resolution.
// This option disables version-related processing during tool resolution.
func WithoutVersion() ResolveOption {
	return func(o *resolveOptions) {
		o.skipVersion = true
	}
}

// WithoutURL returns a ResolveOption that skips URL resolution.
// This option disables URL-related processing during tool resolution.
func WithoutURL() ResolveOption {
	return func(o *resolveOptions) {
		o.skipURL = true
	}
}

// WithUpUntilVersion returns a ResolveOption that enables version checking up to a specific version.
// This option enables version comparison during tool resolution.
func WithUpUntilVersion() ResolveOption {
	return func(o *resolveOptions) {
		o.upUntilVersion = true
	}
}
