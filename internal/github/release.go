package github

type Release struct {
	Name   string `json:"name"`
	Tag    string `json:"tag_name"`
	Assets Assets `json:"assets"`
}
