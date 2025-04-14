//go:generate go tool enumer -type=Extension -output extension_enumer___generated.go -transform=lower
package file

// Extension represents common file extension types.
// Used to categorize files based on their extensions
// and handle platform-specific behaviors.
type Extension int

const (
	// None indicates a file has no extension.
	None Extension = iota

	// EXE represents Windows executable files (.exe).
	EXE

	// GZ represents gzip compressed files (.gz).
	GZ

	// ZIP represents ZIP archive files (.zip).
	ZIP

	// TAR represents tape archive files (.tar).
	TAR

	// Other represents any unrecognized extension.
	Other
)
