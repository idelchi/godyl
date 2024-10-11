package main

import (
	"fmt"

	"github.com/idelchi/godyl/pkg/pretty"
)

type GitHub struct {
	Repo  string
	Owner string
	Token string `mask:"filled"`
}

func main() {
	c := GitHub{
		Repo:  "godyl",
		Owner: "idelchi",
		Token: "something",
	}

	fmt.Println(pretty.JSONMasked(c))
}
