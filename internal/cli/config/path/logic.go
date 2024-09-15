package path

import (
	"fmt"

	"github.com/idelchi/godyl/pkg/path/file"
)

// run executes the `config path` command.
func run(file file.File) error {
	if !file.Exists() {
		return fmt.Errorf("config file %q doesn't exist", file)
	}

	fmt.Println(file)

	return nil
}
