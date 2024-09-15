package platform

// Extension represents a platform-specific file extension.
// Used primarily for executable files and archive formats.
type Extension string

// Default returns the platform's standard executable extension.
// Returns ".exe" for Windows systems and empty string for Unix-like systems.
func (e *Extension) ParseFrom(os OS) {
	switch os.Type() {
	case "windows":
		*e = Extension(".exe")
	default:
		*e = Extension("")
	}
}

// String returns the extension value including the leading dot.
func (e Extension) String() string {
	return string(e)
}

func (e Extension) IsNil() bool {
	return e == ""
}
