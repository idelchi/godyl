package config

// Update holds the configuration options for self-updating the tool.
// These are used as flags, environment variables for the corresponding CLI commands.
type Update struct {
	// Tokens for authentication
	Tokens Tokens `mapstructure:",squash"`

	// Skip SSL verification
	NoVerifySSL bool `mapstructure:"no-verify-ssl"`

	// Version of the tool to install
	Version string

	// Enable pre-release versions
	Pre bool

	// Viper instance
	viperable `mapstructure:"-" yaml:"-" json:"-"`
}
