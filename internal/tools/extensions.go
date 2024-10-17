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
	var noExtensionPart string
	var extensionParts []string

	for _, ext := range exts {
		if ext == "" {
			noExtensionPart = "^[^.]+$"
		} else {
			escapedExt := strings.ReplaceAll(ext, ".", `\.`) // Escape dots in extensions
			extensionParts = append(extensionParts, fmt.Sprintf(".*%s$", escapedExt))
		}
	}

	var pattern string
	// Combine both parts
	if noExtensionPart != "" && len(extensionParts) > 0 {
		pattern = fmt.Sprintf("(%s|%s)", noExtensionPart, strings.Join(extensionParts, "|"))
	} else if noExtensionPart != "" {
		pattern = noExtensionPart
	} else {
		pattern = strings.Join(extensionParts, "|")
	}

	return match.Hint{
		Pattern: pattern,
		Must:    true,
	}
}
