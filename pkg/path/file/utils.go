package file

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/gabriel-vasile/mimetype"
)

// LargerThan checks if the file is larger than the specified size in bytes.
func (f File) LargerThan(size int64) (bool, error) {
	currentSize, err := f.Size()
	if err != nil {
		return false, fmt.Errorf("checking if file %q is larger than %d bytes: %w", f, size, err)
	}

	return currentSize > size, nil
}

// SmallerThan checks if the file is smaller than the specified size in bytes.
func (f File) SmallerThan(size int64) (bool, error) {
	currentSize, err := f.Size()
	if err != nil {
		return false, fmt.Errorf("checking if file %q is smaller than %d bytes: %w", f, size, err)
	}

	return currentSize < size, nil
}

// Hash computes the hash of the file's contents.
func (f File) Hash() (string, error) {
	data, err := f.Read()
	if err != nil {
		return "", fmt.Errorf("reading file %q: %w", f, err)
	}

	hash := sha256.Sum256(data)

	return fmt.Sprintf("%x", hash), nil
}

// IsBinaryLike checks if the file is binary-like.
func (f File) IsBinaryLike() bool {
	if !f.Exists() || !f.IsFile() {
		return false
	}

	// MIME allow-list
	var textish bool
	if m, err := mimetype.DetectFile(f.Path()); err == nil {
		for p := m; p != nil; p = p.Parent() {
			s := p.String()

			if strings.HasPrefix(s, "text/") ||
				p.Is("application/json") ||
				p.Is("application/xml") ||
				p.Is("image/svg+xml") ||
				p.Is("application/javascript") ||
				p.Is("application/yaml") ||
				p.Is("application/toml") ||
				strings.HasSuffix(s, "+json") ||
				strings.HasSuffix(s, "+xml") {
				textish = true
				break
			}
		}
	}

	// Always run content sanity check
	r, err := f.Open()
	if err == nil {
		defer r.Close()
		buf := make([]byte, 128<<10) // read up to 128 KiB
		n, _ := r.Read(buf)
		b := buf[:n]

		// UTF-16/32 BOM detection
		if hasUTF16BOM(b) || hasUTF32BOM(b) {
			// treat as textish unless NUL pattern suggests otherwise
			if bytes.Contains(b, []byte{0}) && !looksLikeUTF16(b) {
				return true
			}
			if textish {
				return false
			}
			return false
		}

		// NUL byte → binary
		if bytes.Contains(b, []byte{0}) {
			return true
		}

		// Invalid UTF-8 → binary
		if !utf8.Valid(b) {
			return true
		}
	}

	// Decide
	if textish {
		return false
	}
	return true
}

// --- Helpers ---
func hasUTF16BOM(b []byte) bool {
	return len(b) >= 2 && ((b[0] == 0xFE && b[1] == 0xFF) || (b[0] == 0xFF && b[1] == 0xFE))
}

func hasUTF32BOM(b []byte) bool {
	return len(b) >= 4 && ((b[0] == 0x00 && b[1] == 0x00 && b[2] == 0xFE && b[3] == 0xFF) ||
		(b[0] == 0xFF && b[1] == 0xFE && b[2] == 0x00 && b[3] == 0x00))
}

func looksLikeUTF16(b []byte) bool {
	// crude check: every other byte is NUL in first few KB
	limit := len(b)
	if limit > 4096 {
		limit = 4096
	}
	nullCount := 0
	for i := 0; i+1 < limit; i += 2 {
		if b[i] == 0x00 || b[i+1] == 0x00 {
			nullCount++
		}
	}
	return nullCount > limit/4 // heuristic threshold
}
