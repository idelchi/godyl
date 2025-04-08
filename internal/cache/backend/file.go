package backend

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/idelchi/godyl/internal/cache/cache"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/folder"
)

// File is a cache backend that stores data in a JSON file
type File struct {
	file.File
	mu    sync.RWMutex
	items cache.Items
}

// NewFile creates a new file backend
func NewFile(path file.File) (*File, error) {
	fb := &File{
		File:  path,
		items: make(cache.Items),
	}

	// Load existing cache data if the file exists
	if fb.Exists() {
		if err := fb.load(); err != nil {
			return nil, err
		}
	} else {
		if err := folder.FromFile(fb.File).Create(); err != nil {
			return nil, err
		}
		if err := fb.Create(); err != nil {
			return nil, err
		} else {
			if err := fb.Write([]byte("[]")); err != nil {
				return nil, err
			}
		}
	}

	return fb, nil
}

// Get retrieves an item from the cache by name
func (fb *File) Get(name string) (*cache.Item, error) {
	fb.mu.RLock()
	defer fb.mu.RUnlock()

	item, ok := fb.items[name]
	if !ok {
		return nil, fmt.Errorf("%w: %q", cache.ErrItemNotFound, name)
	}
	return item, nil
}

// GetAll retrieves all items from the cache
func (fb *File) GetAll() ([]*cache.Item, error) {
	fb.mu.RLock()
	defer fb.mu.RUnlock()

	items := make([]*cache.Item, 0, len(fb.items))
	for _, item := range fb.items {
		items = append(items, item)
	}
	return items, nil
}

// Save stores an item in the cache
func (fb *File) Save(item *cache.Item) error {
	fb.mu.Lock()
	defer fb.mu.Unlock()

	fb.items[item.Name] = item

	return fb.persist()
}

// Delete removes an item from the cache by name
func (fb *File) Delete(name string) error {
	fb.mu.Lock()
	defer fb.mu.Unlock()

	if _, ok := fb.items[name]; !ok {
		return errors.New("item not found")
	}

	delete(fb.items, name)
	return fb.persist()
}

// Close releases any resources held by the backend
func (fb *File) Close() error {
	return nil // No need to close file handles as they're closed after each operation
}

// load reads the cache data from disk
func (fb *File) load() error {
	file, err := fb.Open() // Using embedded File.Open() method
	if err != nil {
		return err
	}
	defer file.Close()

	var items []*cache.Item
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&items); err != nil {
		return err
	}

	fb.items = make(map[string]*cache.Item, len(items))
	for _, item := range items {
		fb.items[item.Name] = item
	}

	return nil
}

// persist writes the cache data to disk
func (fb *File) persist() error {
	file, err := fb.OpenForWriting() // Using embedded File.OpenForWriting() method
	if err != nil {
		return err
	}
	defer file.Close()

	items := make([]*cache.Item, 0, len(fb.items))
	for _, item := range fb.items {
		items = append(items, item)
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty formatting
	return encoder.Encode(items)
}
