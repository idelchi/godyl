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

func FromSlice(slice ...string) Env {
	e := make(Env, len(slice))

	for _, v := range slice {
		kv := strings.SplitN(v, "=", 2)
		e[kv[0]] = kv[1]
	}

	return e
}

func FromEnv() Env {
	return FromSlice(os.Environ()...)
}

func (e *Env) Merge(env Env) {
	for k, v := range env {
		(*e)[k] = v
	}
}
