// Package sources provides abstractions for handling various types of installation sources,
// including GitHub repositories, direct URLs, Go projects, and command-based sources.
// The package defines a common interface, Populater, which is implemented by these sources
// to handle initialization, execution, versioning, path setup, and installation processes.
package tools

type Populator interface {
	Version(*Tool) string
	Download(*Tool) string
}
