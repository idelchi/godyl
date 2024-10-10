package main

import (
	"fmt"
	"os"

	"github.com/inconshreveable/go-update"
)

func doUpdate(file string) error {
	body, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer body.Close()
	if err := update.Apply(body, update.Options{}); err != nil {
		return err
	}
	return err
}
