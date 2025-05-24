package iutils

import "strings"

// Prefix represents a prefix string used for configuration keys.
type Prefix string

// String returns the string representation of the prefix.
func (p Prefix) String() string {
	return string(p)
}

// Lower returns the lowercase representation of the prefix.
func (p Prefix) Lower() Prefix {
	return Prefix(strings.ToLower(p.String()))
}

// RemovePrefix removes the specified prefix from the prefix string.
func (p Prefix) RemovePrefix(prefix string) Prefix {
	return Prefix(strings.TrimPrefix(p.String(), prefix))
}

// RemoveSuffix removes the specified suffix from the prefix string.
func (p Prefix) RemoveSuffix(suffix string) Prefix {
	return Prefix(strings.TrimSuffix(p.String(), suffix))
}

// Upper returns the uppercase representation of the prefix.
func (p Prefix) Upper() Prefix {
	return Prefix(strings.ToUpper(p.String()))
}

// WithUnderscores replaces dots with underscores in the prefix string.
func (p Prefix) WithUnderscores() Prefix {
	return Prefix(strings.ReplaceAll(p.String(), ".", "_"))
}

// Scoped returns a new prefix with an underscore appended.
func (p Prefix) Scoped() Prefix {
	return Prefix(p.String() + "_")
}

// MatchEnvToFlag allows the env provider of koanf to property match stuff like:
// APP_REGISTRY_TOKEN -> registry-token
// when name is for example "APP_".
func MatchEnvToFlag(name Prefix) func(string) string {
	name = name.Lower()

	return func(s string) string {
		s = strings.ToLower(s)
		s = strings.TrimPrefix(s, name.String())

		replacer := strings.NewReplacer(
			"_", "-",
			".", "-",
		)

		// Convert to hyphenated style to match config
		return replacer.Replace(s)
	}
}
