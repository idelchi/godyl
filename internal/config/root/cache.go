package root

import (
	"github.com/idelchi/godyl/pkg/path/folder"
)

type Cache struct {
	// Path to cache folder
	Dir folder.Folder `json:"cache-dir" mapstructure:"cache-dir"`

	// Disabled disables cache interaction
	Disabled bool `json:"no-cache" mapstructure:"no-cache"`
}
