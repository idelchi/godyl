//go:generate go run golang.org/x/tools/cmd/stringer -type Status -output status_string___generated.go
package tools

type Status int

const (
	OK                Status = 200
	BadRequest        Status = 400
	InternalServerErr Status = 500
)
