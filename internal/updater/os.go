package updater

import (
	"runtime"
)

// IsWindows checks if the current operating system is Windows.
// Returns true for Windows systems, false otherwise.
func IsWindows() bool {
	return runtime.GOOS == "windows"
}
