package env

import (
	"fmt"
	"os"
	"strings"
)

type Env map[string]string

func (e Env) Get(key string) (string, error) {
	if v, ok := e[key]; ok {
		return v, nil
	}

	return "", fmt.Errorf("key %q not found", key)
}

func FromEnv() Env {
	env := os.Environ()

	e := make(Env, len(env))

	for _, v := range env {
		kv := strings.SplitN(v, "=", 2)
		e[kv[0]] = kv[1]
	}

	return e
}

func (e *Env) Merge(env Env) {
	for k, v := range env {
		(*e)[k] = v
	}
}
