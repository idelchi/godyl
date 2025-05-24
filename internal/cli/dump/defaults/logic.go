package defaults

import (
	"fmt"

	"github.com/goccy/go-yaml"

	"github.com/idelchi/godyl/internal/cli/common"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/iutils"
)

// run executes the `cache dump` command.
func run(global config.Config, embedded common.Embedded, args ...string) error {
	c, err := getDefaults(embedded, args)
	if err != nil {
		return err
	}

	iutils.Print(global.Dump.Format, c)

	return nil
}

// getDefaults loads and returns the application's default settings.
func getDefaults(files common.Embedded, defaultNames []string) (any, error) {
	var defaults map[string]any

	err := yaml.Unmarshal(files.Defaults, &defaults)
	if err != nil {
		return nil, err
	}

	// If defaultNames is provided, filter the defaults
	if len(defaultNames) > 0 {
		filteredDefaults := make(map[string]any)

		for _, name := range defaultNames {
			if value, exists := defaults[name]; exists {
				filteredDefaults[name] = value
			} else {
				return nil, fmt.Errorf("default %q not found", name)
			}
		}

		defaults = filteredDefaults
	}

	return defaults, nil
}
