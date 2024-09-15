package config

type Status struct {
	Tags []string

	trackable `json:"-" mapstructure:"-" yaml:"-"`
}

func (s Status) ToCommon() Common {
	return Common{
		// Output is not valid for install
		// Output:
		// Strategy is not valid for install
		// Strategy:
		// Source is not valid for install
		// Source:
		// OS is not valid for install
		// OS:
		// Arch is not valid for install
		// Arch:
		// Hints is not valid for install
		// Hints:

		trackable: s.trackable,
	}
}
