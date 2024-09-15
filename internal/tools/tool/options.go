// Package tool provides core functionality for managing tool configurations.
package tool

// ResolveOption is a functional option type for the Resolve method.
type ResolveOption func(*resolveOptions)

// resolveOptions holds all configurable options for the Resolve method.
type resolveOptions struct {
	skipVersion    bool
	upUntilVersion bool
	skipURL        bool
}

func WithoutVersion() ResolveOption {
	return func(o *resolveOptions) {
		o.skipVersion = true
	}
}

func WithoutURL() ResolveOption {
	return func(o *resolveOptions) {
		o.skipURL = true
	}
}

func WithUpUntilVersion() ResolveOption {
	return func(o *resolveOptions) {
		o.upUntilVersion = true
	}
}
