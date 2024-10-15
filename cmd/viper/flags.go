package main

// Config holds all the configuration options for godyl.
type Flags struct {
	Config string
	DotEnv string `mapstructure:"dot-env"`
}
