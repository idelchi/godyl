package tools

import (
	"errors"
	"os"

	"github.com/idelchi/go-next-tag/pkg/stdin"

	"gopkg.in/yaml.v3"
)

// Tools represents a collection of Tool configurations.
type Tools []Tool

// Load reads a tool configuration file and loads it into the Tools collection.
// If the path is "-", it reads from stdin.
// Else, it reads from the specified file path.
func (t *Tools) Load(path string) (err error) {
	var data []byte

	if path == "-" {
		if !stdin.IsPiped() {
			return errors.New("no data piped to stdin")
		}

		input, err := stdin.Read()
		if err != nil {
			return err
		}

		data = []byte(input)
	} else {
		// Read the YAML configuration file from disk.
		input, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		data = input
	}

	// Unmarshal the YAML content into the Tools collection.
	err = yaml.Unmarshal(data, t)
	if err != nil {
		return err
	}

	return nil
}
