package updater

import (
	"runtime"
)

// Operating system identifiers used for platform detection.
const (
	Windows = "windows"
	Linux   = "linux"
	MacOS   = "darwin"
)

// IsWindows checks if the current operating system is Windows.
// Returns true for Windows systems, false otherwise.
func IsWindows() bool {
	return runtime.GOOS == Windows
}

// IsLinux checks if the current operating system is Linux.
// Returns true for Linux systems, false otherwise.
func IsLinux() bool {
	return runtime.GOOS == Linux
}

// IsMacOS checks if the current operating system is macOS.
// Returns true for macOS systems, false otherwise.
func IsMacOS() bool {
	return runtime.GOOS == MacOS
}
