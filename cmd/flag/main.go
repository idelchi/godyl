package main

import (
	"fmt"
	"os"

	"github.com/idelchi/godyl/pkg/pretty"
)

func main() {
	cfg, err := parseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Error validating configuration: %v\n", err)
		os.Exit(1)
	}


	fmt.Println(pretty.JSON(cfg))
}
