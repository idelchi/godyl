package main

import (
	"fmt"
	"os"
)

func main() {
	cfg, err := parseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("config: %v\n", IsSet("config"))
	fmt.Printf(".env: %v\n", IsSet("dot-env"))

	fmt.Printf("config: %v\n", cfg.Config)
	fmt.Printf(".env: %v\n", cfg.DotEnv)
}
