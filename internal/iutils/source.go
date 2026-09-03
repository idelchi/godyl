package iutils

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/stdin"
)

// Source represents a source of input data.
type Source interface {
	// Read returns all bytes available from the source.
	Read() ([]byte, error)
}

// FileSource reads from a file.
type FileSource struct {
	// File is the path read by Read.
	File file.File
}

// Read returns the complete contents of the configured file.
func (s FileSource) Read() ([]byte, error) {
	content, err := s.File.Read()
	if err != nil {
		return nil, err
	}

	return content, nil
}

// StdinSource reads from stdin.
type StdinSource struct{}

// Read returns all data currently piped through standard input.
func (s StdinSource) Read() ([]byte, error) {
	input, err := stdin.Read()
	if err != nil {
		return nil, fmt.Errorf("reading from stdin: %w", err)
	}

	return []byte(input), nil
}

// BytesSource represents pre-loaded bytes.
type BytesSource struct {
	// Data contains the bytes returned by Read.
	Data []byte
}

// Read returns the stored byte slice without copying it.
func (s BytesSource) Read() ([]byte, error) {
	return s.Data, nil
}

// MultiSource reads from multiple sources and concatenates the results.
type MultiSource struct {
	// Sources are read in order and separated with a newline.
	Sources []Source
}

// NewMultiSource creates a new multi source.
func NewMultiSource(sources ...Source) *MultiSource {
	return &MultiSource{Sources: sources}
}

// Read concatenates sources in order with one newline between each result.
func (s MultiSource) Read() ([]byte, error) {
	var buf bytes.Buffer

	for i, source := range s.Sources {
		data, err := source.Read()
		if err != nil {
			return nil, fmt.Errorf("reading from source: %w", err)
		}

		if i > 0 {
			buf.WriteByte('\n')
		}

		buf.Write(data)
	}

	return buf.Bytes(), nil
}

// ErrInvalidSource is returned when the source is invalid.
var ErrInvalidSource = errors.New("invalid source")

// GetSourceFromPath determines the appropriate source type for a given path.
func GetSourceFromPath(path string) (Source, error) {
	if path == "-" {
		if !stdin.IsPiped() {
			return nil, fmt.Errorf("%w: '-' indicated stdin, but no data is piped", ErrInvalidSource)
		}

		return StdinSource{}, nil
	}

	return FileSource{File: file.New(path)}, nil
}

// ReadPaths reads from one or more file paths or "-" (stdin), and returns concatenated data.
func ReadPaths(paths ...string) ([]byte, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("%w: no paths provided", ErrInvalidSource)
	}

	sources := make([]Source, 0, len(paths))

	for _, path := range paths {
		source, err := GetSourceFromPath(path)
		if err != nil {
			return nil, fmt.Errorf("getting source from path %q: %w", path, err)
		}

		sources = append(sources, source)
	}

	if len(sources) == 1 {
		return sources[0].Read() //nolint:wrapcheck // Error does not need additional wrapping.
	}

	return NewMultiSource(sources...).Read()
}

// ReadPathsOrDefault reads from args if provided, or from defaultPath if args is empty.
func ReadPathsOrDefault(defaultPath string, args ...string) ([]byte, error) {
	if len(args) == 0 {
		return ReadPaths(defaultPath)
	}

	return ReadPaths(args...)
}
