package config

import (
	"fmt"
	"slices"

	"github.com/idelchi/godyl/internal/tools"
)

// ValidateUpdateStrategy checks if the update strategy is valid.
func ValidateUpdateStrategy(strategy tools.Strategy) error {
	allowedUpdateStrategies := []tools.Strategy{tools.None, tools.Upgrade, tools.Force}
	if !slices.Contains(allowedUpdateStrategies, strategy) {
		return fmt.Errorf(
			"%w: unknown update strategy: %q: allowed are %v",
			ErrUsage,
			strategy,
			allowedUpdateStrategies,
		)
	}
	return nil
}

// ValidateSourceType checks if the source type is valid.
func ValidateSourceType(sourceType string) error {
	allowedSourceTypes := []string{"github", "url", "go", "command"}
	if !slices.Contains(allowedSourceTypes, sourceType) {
		return fmt.Errorf(
			"%w: unknown source type: %q: allowed are %v",
			ErrUsage,
			sourceType,
			allowedSourceTypes,
		)
	}
	return nil
}
