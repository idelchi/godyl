package config

// Config holds the top level configuration for godyl.
// It is split into sub-structs for each command.

type Config struct {
	// Root level configuration, mapping configurations on the root `godyl` command
	Root Root

	// Tool level configuration, mapping configurations on the `install`, `download`,
	// and (partially) the `update` commands
	Tool Tool

	// Dump level configuration, mapping configurations on the `dump` command
	Dump Dump
}
