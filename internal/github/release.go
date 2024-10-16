package github

// Release represents a GitHub release, containing the release name, tag, and associated assets.
type Release struct {
	Name   string `json:"name"`     // Name is the name of the release.
	Tag    string `json:"tag_name"` // Tag is the tag associated with the release (e.g., version number).
	Assets Assets `json:"assets"`   // Assets is a collection of assets attached to the release.
}
