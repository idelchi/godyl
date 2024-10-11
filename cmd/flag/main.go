package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-envparse"

	"mvdan.cc/sh/v3/expand"
)

type Env map[string]string

func (e Env) Get(key string) string {
	return e[key]
}

func (e *Env) Merge(env Env) {
	for k, v := range env {
		(*e)[k] = v
	}
}

func Get(t *ast.Task) []string {
	if t.Env == nil {
		return nil
	}
	environ := os.Environ()
	for k, v := range t.Env.ToCacheMap() {
		if !isTypeAllowed(v) {
			continue
		}
		if !experiments.EnvPrecedence.Enabled {
			if _, alreadySet := os.LookupEnv(k); alreadySet {
				continue
			}
		}
		environ = append(environ, fmt.Sprintf("%s=%v", k, v))
	}

	return environ
}

func main() {
	file, err := os.Open("cmd/flag/.env")
	if err != nil {
		panic(err)
	}
	env, err := envparse.Parse(file)
	if err != nil {
		panic(err)
	}

	Env := expand.ListEnviron(os.Environ()...)

	fmt.Println(env)

	fmt.Println(Env.Get("ALLUSERSPROFILE"))
}
