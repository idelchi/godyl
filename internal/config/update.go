package config

import (
	"github.com/idelchi/godyl/internal/tools"
)

type Update struct {
	// Strategy to use for updating tools
	Strategy tools.Strategy `validate:"upgrade force"`

	// Tokens for authentication
	Tokens Tokens `mapstructure:",squash"`

	// Skip SSL verification
	NoVerifySSL bool `mapstructure:"no-verify-ssl"`
}
