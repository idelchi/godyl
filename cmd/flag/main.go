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
	// Test FromEnv
	envVars := env.FromEnv()
	// fmt.Println("Original Environment Variables:", envVars.ToSlice())

	// Test Normalized for Windows (won't modify on non-Windows)
	// normalizedEnv := envVars.Normalized()
	// fmt.Println("Normalized Environment Variables:", normalizedEnv.ToSlice())

	// Test Get
	key := "PATH"
	value, err := envVars.Get(key)
	if err != nil {
		fmt.Printf("Error getting key %q: %v\n", key, err)
	} else {
		fmt.Printf("Value of key %q: %s\n", key, value)
	}

	// Test ToSlice
	// envSlice := envVars.ToSlice()
	// fmt.Println("Environment variables as slice:", envSlice)

	// Test FromSlice
	newEnvVars := env.FromSlice("KEY1=VALUE1=SHART", "KEY2=VALUE2")
	fmt.Println("New Env from Slice:", newEnvVars.ToSlice())

	// Test Add
	newEnvVars.Add("KEY3=VALUE3=SHIT")
	fmt.Println("Env after Add:", newEnvVars.ToSlice())

	key = "KEY1"
	value, err = newEnvVars.Get(key)
	if err != nil {
		fmt.Printf("Error getting key %q: %v\n", key, err)
	} else {
		fmt.Printf("Value of key %q: %s\n", key, value)
	}

	pretty.PrintJSON(newEnvVars)

	// Test Merge
	// envVars.Merge(newEnvVars)
	// fmt.Println("Env after Merge:", envVars.ToSlice())

	// Setting a custom environment variable for testing
	os.Setenv("CUSTOM_KEY", "CUSTOM_VALUE")
	// envVars = env.FromEnv()
	// fmt.Println("Updated Environment Variables:", envVars.ToSlice())
}
