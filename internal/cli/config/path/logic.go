package path

import (
	"fmt"

	"github.com/idelchi/godyl/pkg/path/file"
)

// run executes the `config path` command.
func run(file file.File) error {
	fmt.Println(file)

	return nil
}
