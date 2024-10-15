package goi

type Release struct {
	Version string   `json:"version"`
	Files   []Target `json:"files"`
}
