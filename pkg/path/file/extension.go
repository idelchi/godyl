//go:generate go tool enumer -type=Extension -output extension_enumer___generated.go
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
	// TAR represents the ".tar" file extension.
	TAR
	// Other represents any other, unrecognized file extension.
	Other
)
