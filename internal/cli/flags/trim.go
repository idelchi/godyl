package flags

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/idelchi/godyl/pkg/path/file"
	"gopkg.in/yaml.v3"
)

// Trim reads a YAML file and returns an io.Reader with only the content
// under the specified prefix, with the prefix itself removed.
func Trim(file file.File, prefix string) (io.Reader, error) {
	// Read the YAML file
	data, err := file.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	if prefix == "" {
		return bytes.NewReader(data), nil
	}

	// Parse the YAML into a map
	var yamlData map[string]any
	if err := yaml.Unmarshal(data, &yamlData); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Navigate through the nested structure based on the prefix
	prefixParts := strings.Split(prefix, ".")
	result, found := extractNestedData(yamlData, prefixParts)
	if !found {
		// If not found, return an empty map
		return bytes.NewReader([]byte("{}")), nil
	}

	// Marshal the result back to YAML
	resultYAML, err := yaml.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result to YAML: %w", err)
	}

	// Return the result as an io.Reader
	return bytes.NewReader(resultYAML), nil
}

// extractNestedData navigates the YAML map to find the data at the specified path
func extractNestedData(data map[string]any, path []string) (any, bool) {
	if len(path) == 0 {
		return nil, false
	}

	current := path[0]
	if value, exists := data[current]; exists {
		if len(path) == 1 {
			// We've reached the target prefix
			return value, true
		}

		// Continue navigating if we have more path segments and the current value is a map
		if nestedMap, ok := value.(map[string]any); ok {
			return extractNestedData(nestedMap, path[1:])
		}
	}

	return nil, false
}
