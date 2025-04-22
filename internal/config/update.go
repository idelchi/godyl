package config

// Update holds the configuration options for self-updating the tool.
// These are used as flags, environment variables for the corresponding CLI commands.
type Update struct {
	viperable   `json:"-" mapstructure:"-" yaml:"-"`
	Tokens      Tokens `mapstructure:",squash"`
	Version     string
	NoVerifySSL bool `mapstructure:"no-verify-ssl"`
	Pre         bool
	Check       bool
	Cleanup     bool
}
