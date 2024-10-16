package file

// Extension represents a file extension.
type Extension int

const (
	// None represents no file extension.
	None Extension = iota
	// EXE represents the ".exe" file extension.
	EXE
	// GZ represents the ".gz" file extension.
	GZ
	// ZIP represents the ".zip" file extension.
	ZIP
	// Other represents any other file extension.
	Other
)
