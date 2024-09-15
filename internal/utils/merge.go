package utils

import (
	"fmt"

	"dario.cat/mergo"
	"github.com/idelchi/godyl/pkg/utils"
)

func Merge[T any](dst *T, src *T, opts ...func(*mergo.Config)) (err error) {
	if len(opts) == 0 {
		// Default options for merging.
		// TODO(Idelchi): Will `mergo.WithoutDereference` cause issues for some edge cases?
		opts = []func(*mergo.Config){mergo.WithOverride, mergo.WithoutDereference}
	}

	// Avoid pointers being copied. As such, we can always "copy" both the source and destination.
	src, err = utils.DeepCopyPtr(src)
	if err != nil {
		return fmt.Errorf("copying source in preparation for merge: %w", err)
	}

	// dst, err = utils.DeepCopyPtr(dst)
	// if err != nil {
	// 	return fmt.Errorf("copying destination in preparation for merge: %w", err)
	// }

	err = mergo.Merge(dst, src, opts...)

	return err
}
