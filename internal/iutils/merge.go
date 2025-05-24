package iutils

import (
	"fmt"

	"github.com/idelchi/godyl/pkg/utils"

	"dario.cat/mergo"
)

// Merge merges the source and destination values.
// It uses the mergo package to perform the merge.
// The destination value is modified in place.
// The source value is deep-copied to avoid modifying the original.
// The function accepts options for the merge operation.
// If no options are provided, it uses the default options:
// mergo.WithOverride and mergo.WithoutDereference.
func Merge[T any](dst, src *T, opts ...func(*mergo.Config)) (err error) {
	if len(opts) == 0 {
		// Default options for merging.
		opts = []func(*mergo.Config){mergo.WithOverride, mergo.WithoutDereference}
	}

	// Avoid pointers being copied. As such, we can always "copy" both the source and destination.
	src, err = utils.DeepCopyPtr(src)
	if err != nil {
		return fmt.Errorf("copying source in preparation for merge: %w", err)
	}

	err = mergo.Merge(dst, src, opts...)

	return err //nolint:wrapcheck // Error does not need additional wrapping.
}
