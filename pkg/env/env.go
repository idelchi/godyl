package env

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

type Env map[string]string

func (e Env) Normalized() Env {
	if runtime.GOOS == "windows" {
		normalized := make(Env, len(e))

		for k, v := range e {
			normalized[strings.ToUpper(k)] = v
		}

		return normalized
	}

	return e
}

func (e Env) Get(key string) (string, error) {
	if v, ok := e[key]; ok {
		return v, nil
	}

	return "", fmt.Errorf("key %q not found", key)
}

func (e Env) ToSlice() []string {
	slice := make([]string, 0, len(e))

	for k, v := range e {
		slice = append(slice, k+"="+v)
	}

	return slice
}

func FromSlice(slice ...string) Env {
	e := make(Env, len(slice))

	for _, v := range slice {
		kv := strings.SplitN(v, "=", 2)
		e[kv[0]] = kv[1]
	}

	return e
}

func (e *Env) Add(slice ...string) {
	e.Merge(FromSlice(slice...))
}

func FromEnv() Env {
	return FromSlice(os.Environ()...)
}

func (e *Env) Merge(env Env) {
	for k, v := range env {
		if _, ok := (*e)[k]; ok {
			continue
		}

		(*e)[k] = v
	}
}
