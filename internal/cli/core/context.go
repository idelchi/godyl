package core

import (
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/koanfx"
)

// GlobalContext is a global variable, allowing the root parsed configuration to be
// accessed from anywhere in the application.
//
//nolint:gochecknoglobals 	// Necessary for global access to the parsed configuration.
var GlobalContext Context

// Context holds the parsed configuration.
type Context struct {
	// Config holds a map of the parsed configuration file.
	Config *koanfx.Koanf
	// DotEnv holds the parsed .env files
	DotEnv *env.Env
	// Env holds the parsed environment variables.
	Env *env.Env
}
