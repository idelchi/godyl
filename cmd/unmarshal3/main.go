package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/idelchi/godyl/pkg/pretty"
	"github.com/idelchi/godyl/pkg/unmarshal"
	"gopkg.in/yaml.v3"
)

type Person struct {
	Name string `yaml:"name"`
	Age  int    `yaml:"age"`
}

// type (
// 	Age  int
// 	Name string
// )

// Type aliases using the generic wrapper
type (
	Ages  = unmarshal.SingleOrSlice[int]
	Names = unmarshal.SingleOrSlice[string]
)

type Persons []Person

func (p *Persons) UnmarshalYAML(value *yaml.Node) error {
	result, err := unmarshal.UnmarshalSingleOrSlice[Person](value, true)
	if err != nil {
		return err
	}
	*p = result
	return nil
}

func (p Persons) Oldest() int {
	max := 0
	for _, person := range p {
		if person.Age > max {
			max = person.Age
		}
	}
	return max
}

type Config struct {
	Persons Persons `yaml:"persons"`
	Ages    Ages    `yaml:"ages"`
	Names   Names   `yaml:"names"`
}

func main() {
	cfg := flag.String("cfg", "single", "config file")

	flag.Parse()

	*cfg = "cmd/unmarshal3/config-" + *cfg + ".yml"

	file, err := os.ReadFile(*cfg)
	if err != nil {
		fmt.Println(err)

		os.Exit(1)
	}

	// Unmarshal the content into a map for pretty printing
	var content map[string]any
	err = yaml.Unmarshal(file, &content)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("File content:")
	pretty.PrintJSON(content)

	_, err = os.Open(*cfg)
	if err != nil {
		fmt.Println(err)

		os.Exit(1)
	}
	var config Config

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Config:")
	pretty.PrintJSON(config)

	// Unmarshal the file content into a Config struct
	// var config Config
	// dec := yaml.NewDecoder(filex)
	// dec.KnownFields(true)
	// if err := dec.Decode(&config); err != nil {
	// 	fmt.Println(err)

	// 	os.Exit(1)
	// }
}
