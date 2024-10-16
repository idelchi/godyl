package file

// Extension represents a file extension type.
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
	// Other represents any other, unrecognized file extension.
	Other
)
