package root

// Tokens holds the configuration options for authentication tokens.
type Tokens struct {
	// GitHub token for authentication
	GitHub string `json:"github-token" mapstructure:"github-token" mask:"fixed"`
	// GitLab token for authentication
	GitLab string `json:"gitlab-token" mapstructure:"gitlab-token" mask:"fixed"`
	// URL token for authentication
	URL string `json:"url-token" mapstructure:"url-token" mask:"fixed"`
}
