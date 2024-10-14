package match

type Hints []Hint

func (h *Hints) Add(hints Hints) {
	for _, hint := range hints {
		*h = append(*h, hint)
	}
}

type Hint struct {
	Pattern        string
	WeightTemplate string `json:"-" yaml:"weight" mapstructure:"weight"`
	Weight         int    `yaml:"-" mapstructure:"-"`
	Regex          bool
	Must           bool
}

func NewDefaultHint(pattern string) Hint {
	return Hint{
		Pattern: pattern,
		Weight:  1,
		Must:    true,
	}
}
