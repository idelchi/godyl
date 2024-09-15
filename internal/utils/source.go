package utils

import (
	"errors"

	"github.com/idelchi/go-next-tag/pkg/stdin"
	"github.com/idelchi/godyl/pkg/path/file"
)

// Source represents a source of input data
type Source interface {
	// Read returns the data from the source
	Read() ([]byte, error)
}

// FileSource reads from a file
type FileSource struct {
	file file.File
}

// NewFileSource creates a new file source
func NewFileSource(file file.File) FileSource {
	return FileSource{file: file}
}

// Read implements Source
func (s FileSource) Read() ([]byte, error) {
	return s.file.Read()
}

// StdinSource reads from stdin
type StdinSource struct{}

// NewStdinSource creates a new stdin source
func NewStdinSource() *StdinSource {
	return &StdinSource{}
}

// Read implements Source
func (s StdinSource) Read() ([]byte, error) {
	input, err := stdin.Read()
	if err != nil {
		return nil, err
	}

	return []byte(input), nil
}

// BytesSource represents pre-loaded bytes
type BytesSource struct {
	data []byte
}

// NewBytesSource creates a new bytes source
func NewBytesSource(data []byte) *BytesSource {
	return &BytesSource{data: data}
}

// Read implements Source
func (s BytesSource) Read() ([]byte, error) {
	return s.data, nil
}

// MultiSource reads from multiple sources and concatenates the results
type MultiSource struct {
	sources []Source
}

// NewMultiSource creates a new multi source
func NewMultiSource(sources ...Source) *MultiSource {
	return &MultiSource{sources: sources}
}

// Read implements Source
func (s MultiSource) Read() ([]byte, error) {
	var result []byte

	for _, source := range s.sources {
		data, err := source.Read()
		if err != nil {
			return nil, err
		}

		result = append(result, data...)
	}

	return result, nil
}

// GetSourceFromFile determines the appropriate source type for a file
func GetSourceFromFile(file file.File) (Source, error) {
	if file.Path() == "-" {
		if !stdin.IsPiped() {
			return nil, errors.New("no data piped to stdin")
		}
		return NewStdinSource(), nil
	}

	return NewFileSource(file), nil
}

// ReadSingle reads data from a single file or stdin
func ReadSingle(file file.File) ([]byte, error) {
	source, err := GetSourceFromFile(file)
	if err != nil {
		return nil, err
	}

	return source.Read()
}

// ReadMultiple reads data from multiple files and concatenates the results
func ReadMultiple(files []file.File) ([]byte, error) {
	if len(files) == 0 {
		return nil, errors.New("no files provided")
	}

	if len(files) == 1 {
		return ReadSingle(files[0])
	}

	sources := make([]Source, 0, len(files))

	for _, f := range files {
		source, err := GetSourceFromFile(f)
		if err != nil {
			return nil, err
		}
		sources = append(sources, source)
	}

	return NewMultiSource(sources...).Read()
}

// ReadFromArgs reads data from CLI args
// Args can be a list of files, nothing, or "-" for stdin
func ReadFromArgs(defaultTool string, args ...string) ([]byte, error) {
	// If no args provided, use default file in the base path
	if len(args) == 0 {
		defaultFile := file.New(defaultTool)
		return ReadSingle(defaultFile)
	}

	// Check for stdin
	if len(args) == 1 && args[0] == "-" {
		stdinFile := file.New("", "-")
		return ReadSingle(stdinFile)
	}

	// Handle multiple files
	files := make([]file.File, 0, len(args))
	for _, arg := range args {
		files = append(files, file.New(arg))
	}

	return ReadMultiple(files)
}
