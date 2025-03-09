// Package cobraext provides extensions to the cobra package.
// It provides an UnknownSubcommandAction function that prints an error message for unknown subcommands.
// This is necessary when using `TraverseChildren: true`,
// because it seems to disable suggestions for unknown subcommands.
package cobraext
