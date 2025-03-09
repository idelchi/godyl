package config

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

// ErrExitGracefully is an error that signals the program to exit gracefully.
var ErrExitGracefully = errors.New("gracefully exiting")

// Validate performs validation checks.
func Validate(validations ...any) error {
	for _, v := range validations {
		if err := validator.New().Struct(v); err != nil {
			return fmt.Errorf("validating config: %w\nSee --help for more info on usage", err)
		}
	}

	return nil
}
