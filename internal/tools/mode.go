package tools

type Mode string

const (
	Extract Mode = "extract"
	Find    Mode = "find"
)

func (m Mode) String() string {
	return string(m)
}
