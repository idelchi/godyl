//go:generate go tool enumer -type=Status -output status_enumer___generated.go -transform=lower
package tools

// Result represents the outcome of an installation process with additional context
type Result struct {
	Status  Status
	Message string
}

// ResultStatus represents the possible statuses of an installation process
type Status int

const (
	// Installed means the tool was installed successfully
	Installed Status = iota
	// AlreadyInstalled means the tool was already installed and thus skipped
	AlreadyInstalled
	// Updated means the tool was updated successfully
	Updated
	// Skipped means the tool was skipped for other reasons
	Skipped
	// Failed means the installation process failed
	Failed
	// Resolved means the tool was resolved successfully
	Resolved
)
