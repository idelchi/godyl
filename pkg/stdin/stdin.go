// Package stdin provides simple utilities for reading from stdin,
// as well as determining if stdin is a terminal or if something has been piped to it.
package stdin

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

// IsInteractive reports whether the program is attached to a real terminal.
func IsInteractive() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}

// HasInput reports true when stdin is **not** a terminal.
// This covers pipes (echo foo | app), file redirects (app < file),
// and CI/dev‑null situations.  It never blocks.
func HasInput() bool {
	return !IsInteractive()
}

// IsPiped checks if something has been piped to stdin.
func IsPiped() bool {
	fi, err := os.Stdin.Stat()

	return (fi.Mode()&os.ModeCharDevice) == 0 && err == nil
}

// MaybePiped checks if something has been piped to stdin.
func MaybePiped() (bool, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false, fmt.Errorf("getting stdin stat: %w", err)
	}

	isPipe := (stat.Mode()&os.ModeNamedPipe) != 0 ||
		(stat.Mode()&(os.ModeCharDevice|os.ModeDir|os.ModeSymlink)) == 0

	return isPipe, nil
}

// Read returns stdin as a string, trimming the trailing newline.
func Read() (string, error) {
	bytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		err = fmt.Errorf("reading from stdin: %w", err)
	}

	return strings.TrimSuffix(string(bytes), "\n"), err
}

// ReadAll consumes stdin and returns the trimmed string plus a flag that
// tells whether any bytes were actually read.
func ReadAll() (data string, read bool, err error) {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", false, err
	}

	if len(b) == 0 {
		return "", false, nil
	}

	return strings.TrimSuffix(string(b), "\n"), true, nil
}
