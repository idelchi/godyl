package updater

import (
	"runtime"
)

// Platform-specific constants and helper functions
const (
	Windows = "windows"
	Linux   = "linux"
	MacOS   = "darwin"
)

// IsWindows returns true if the current OS is Windows.
func IsWindows() bool {
	return runtime.GOOS == Windows
}

// IsLinux returns true if the current OS is Linux.
func IsLinux() bool {
	return runtime.GOOS == Linux
}

// IsMacOS returns true if the current OS is macOS.
func IsMacOS() bool {
	return runtime.GOOS == MacOS
}
