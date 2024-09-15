//nolint:forbidigo // Functions in file will print & exit for various help messages.
package pretty

import (
	"fmt"
	"strings"
)

// PrintJSON outputs formatted JSON to stdout.
// Prints any value as indented JSON with consistent formatting.
func PrintJSON(obj any) {
	fmt.Println(JSON(obj))
}

// PrintYAML outputs formatted YAML to stdout.
// Prints any value as indented YAML with consistent formatting.
func PrintYAML(obj any) {
	fmt.Println(YAML(obj))
}

// PrintJSONMasked outputs masked JSON to stdout.
// Prints any value as JSON with sensitive fields masked for
// security. Uses MaskJSON internally for field masking.
func PrintJSONMasked(obj any) {
	fmt.Println(strings.TrimSpace(JSONMasked(obj)))
}

// PrintYAMLMasked outputs masked YAML to stdout.
// Prints any value as YAML with sensitive fields masked for
// security. Uses MaskYAML internally for field masking.
func PrintYAMLMasked(obj any) {
	fmt.Println(strings.TrimSpace(YAMLMasked(obj)))
}

// PrintEnv outputs environment variables to stdout.
// Prints any value in KEY=VALUE format suitable for .env files.
func PrintEnv(env any) {
	fmt.Println(strings.TrimSpace(Env(env)))
}
