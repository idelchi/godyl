package pretty

import (
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/joho/godotenv"
	"github.com/showa-93/go-mask"
)

const indent = 2

var YAMLOptions = []yaml.EncodeOption{
	yaml.Indent(indent),                   // Set indentation to 2 spaces
	yaml.UseSingleQuote(true),             // Use single quotes for strings
	yaml.UseLiteralStyleIfMultiline(true), // Use literal style for multiline strings
}

// YAML formats data as indented YAML.
// Converts any value to a formatted YAML string with consistent
// indentation. Returns error message as string if encoding fails.
func YAML(obj any) string {
	// Use MarshalWithOptions to set the indent
	yamlBytes, err := yaml.MarshalWithOptions(
		obj,
		YAMLOptions...,
	)
	if err != nil {
		return err.Error()
	}

	return string(yamlBytes)
}

// YAMLMasked formats data as YAML with sensitive data masked.
// Converts any value to YAML, first applying masking rules to
// hide sensitive fields. Uses JSONMasked internally for masking.
func YAMLMasked(obj any) string {
	return YAML(MaskJSON(obj))
}

// JSON formats data as indented JSON.
// Converts any value to a formatted JSON string with consistent
// indentation. Returns error message as string if marshaling fails.
func JSON(obj any) string {
	bytes, err := json.MarshalIndent(obj, "  ", "    ")
	if err != nil {
		return err.Error()
	}

	return string(bytes)
}

// JSONMasked formats data as JSON with sensitive data masked.
// Converts any value to JSON, first applying masking rules to
// hide sensitive fields like passwords and tokens.
func JSONMasked(obj any) string {
	return JSON(MaskJSON(obj))
}

// MaskJSON applies data masking rules to an object.
// Uses predefined rules to identify and mask sensitive fields,
// replacing them with "-" characters. Returns the masked object
// or error message if masking fails.
func MaskJSON(obj any) any {
	masker := mask.NewMasker()

	// Set the masking character to "-"
	masker.SetMaskChar("-")

	// Register a function to mask strings by filling them with the masking character.
	masker.RegisterMaskStringFunc(mask.MaskTypeFilled, masker.MaskFilledString)

	t, err := mask.Mask(obj)
	if err != nil {
		return err.Error()
	}

	return t
}

// Env formats data as environment variables.
// Converts any value to KEY=VALUE format suitable for .env files.
// Complex objects are flattened to string representations.
func Env(obj any) string {
	// Convert to map via JSON
	jsonData, err := json.Marshal(obj)
	if err != nil {
		return err.Error()
	}

	var data map[string]any
	if err = json.Unmarshal(jsonData, &data); err != nil {
		return err.Error()
	}

	// Convert to string map (godotenv requires map[string]string)
	stringMap := make(map[string]string)
	for k, v := range data {
		// Convert each value to string
		stringMap[k] = fmt.Sprintf("%v", v)
	}

	// Use godotenv.Marshal to format in dotenv style
	result, err := godotenv.Marshal(stringMap)
	if err != nil {
		return err.Error()
	}

	return result
}
