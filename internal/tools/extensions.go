package tools

import (
	"fmt"
	"strings"

	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/pkg/unmarshal"
)

// Extensions represents a collection of file extensions related to the tool.
// It can either be a single string or a slice of strings, providing flexibility
// when configuring tools that may involve multiple file types.
type Extensions = unmarshal.SingleOrSlice[string]

func ExtensionsToHints(exts Extensions) match.Hint {
	var hints Hint
	for _, ext := range exts {
		hints = append(hints, NewDefaultHint(`\.`+ext+`$`))
	}
	return hints
}

func buildRegex(extensions []string) string {
	var noExtensionPart string
	var extensionParts []string

	for _, ext := range extensions {
		if ext == "" {
			noExtensionPart = "^[^.]+$"
		} else {
			escapedExt := strings.ReplaceAll(ext, ".", `\.`) // Escape dots in extensions
			extensionParts = append(extensionParts, fmt.Sprintf(".*%s$", escapedExt))
		}
	}

	// Combine both parts
	if noExtensionPart != "" && len(extensionParts) > 0 {
		return fmt.Sprintf("(%s|%s)", noExtensionPart, strings.Join(extensionParts, "|"))
	} else if noExtensionPart != "" {
		return noExtensionPart
	} else {
		return strings.Join(extensionParts, "|")
	}
}

func main() {
	extensions := []string{"", ".zip", ".tar.gz"}
	regex := buildRegex(extensions)
	fmt.Println(regex)
}
