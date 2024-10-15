package tools

import (
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/pkg/unmarshal"
)

type Fallbacks = unmarshal.SingleOrSlice[sources.Type]
