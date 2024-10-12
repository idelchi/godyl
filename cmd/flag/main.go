package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-envparse"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/pretty"

	"mvdan.cc/sh/v3/expand"
)

func main2() {
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

func main() {
	os.Setenv("CUSTOM_KEY", "CUSTOM_VALUE")
	// Test FromEnv
	envVars := env.FromEnv()

	// Test FromSlice
	newEnvVars := env.FromSlice("KEY1=VALUE1", "KEY2=VALUE2")

	// Test Add
	newEnvVars.Add("KEY3=SKR54/=//SHIT")
	newEnvVars.Add("HOSTNAME=SHIT")

	key := "KEY1"
	value, err := newEnvVars.Get(key)
	if err != nil {
		fmt.Printf("Error getting key %q: %v\n", key, err)
	} else {
		fmt.Printf("Value of key %q: %s\n", key, value)
	}

	fmt.Println("Before Merge")
	pretty.PrintJSON(newEnvVars)

	// Test Merge
	newEnvVars.Merge(envVars)

	fmt.Println("After merge")
	pretty.PrintJSON(newEnvVars)

	// Setting a custom environment variable for testing
	// envVars = env.FromEnv()
	// fmt.Println("Updated Environment Variables:", envVars.ToSlice())
}
