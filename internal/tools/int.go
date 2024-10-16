//go:generate stringer -type Status
package tools

type Status int

const (
	OK                Status = 200
	BadRequest        Status = 400
	InternalServerErr Status = 500
)
