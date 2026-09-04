// Package ai defines configuration for the optional AI asset selector.
package ai

// Config controls the AI tie-breaker used for ambiguous asset matches.
type Config struct {
	// Enabled permits AI selection between equally ranked candidates.
	Enabled bool `mapstructure:"ai" yaml:"ai"`
	// Provider selects the Fantasy provider adapter.
	Provider string `mapstructure:"ai-provider" yaml:"ai-provider"`
	// Model overrides the provider's default model.
	Model string `mapstructure:"ai-model" yaml:"ai-model"`
	// URL overrides the provider's API endpoint.
	URL string `mapstructure:"ai-url" yaml:"ai-url"`
	// APIKey authenticates the selected provider.
	APIKey string `mapstructure:"ai-api-key" mask:"fixed" yaml:"ai-api-key"`
}
