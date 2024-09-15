package utils

import "strings"

type Prefix string

func (p Prefix) String() string {
	return string(p)
}

func (p Prefix) Lower() Prefix {
	return Prefix(strings.ToLower(p.String()))
}

func (p Prefix) RemovePrefix(prefix string) Prefix {
	return Prefix(strings.TrimPrefix(p.String(), prefix))
}

func (p Prefix) Upper() Prefix {
	return Prefix(strings.ToUpper(p.String()))
}

func (p Prefix) WithUnderscores() Prefix {
	return Prefix(strings.ReplaceAll(p.String(), ".", "_"))
}

func (p Prefix) Scoped() Prefix {
	return Prefix(p.String() + "_")
}

// MatchEnvToFlag allows the env provider of koanf to property match stuff like:
// APP_REGISTRY_TOKEN -> registry-token
// when name is for example "APP_"
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
