package goi

// Release represents a Go release, containing the version and a list of files (targets) associated with the release.
type Release struct {
	Version string   `json:"version"` // Version specifies the version of the Go release.
	Files   []Target `json:"files"`   // Files contains the list of files (targets) available for this release.
}
