// Package cobraext provides extensions to the cobra package.
// It provides a NewDefaultRootCommand function that creates a root command with default settings.
// It sets up integration with viper, with environment variable and flag binding.
// Additional functions can be passed to be executed before the command is run.
// It also provides an UnknownSubcommandAction function that prints an error message for unknown subcommands.
// This is necessary when using `TraverseChildren: true`,
// because it seems to disable suggestions for unknown subcommands.
package cobraext
