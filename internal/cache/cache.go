// Package cache returns a default cache manager with the `File` backend.
package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/idelchi/godyl/internal/tools/version"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/folder"
)

// New creates a new cache manager with the specified file as the backend.
func New(file file.File) *Cache {
	return &Cache{
		File:  file,
		items: make(Items),
	}
}

// Item represents a cache item.
type Item struct {
	Version    version.Version `json:"version"`
	Downloaded time.Time       `json:"downloaded"`
	Updated    time.Time       `json:"updated"`
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Path       string          `json:"path"`
	Type       string          `json:"type"`
}

// ErrItemNotFound is returned when an item is not found in the cache.
var ErrItemNotFound = errors.New("item not found")

// Items is a map of items indexed by their names.
type Items map[string]*Item

// Cache is a cache backend that stores data in a JSON file.
type Cache struct {
	items Items
	file.File
	mu sync.RWMutex

	wasTouched bool
}

// Load creates or loads the cache file.
func (c *Cache) Load() error {
	// Load existing cache data if the file exists
	if c.Exists() {
		return c.load()
	}

	err := folder.FromFile(c.File).Create()
	if err != nil {
		return err //nolint:wrapcheck 	// Error does not need additional wrapping.
	}

	err = c.Create()
	if err != nil {
		return err //nolint:wrapcheck 	// Error does not need additional wrapping.
	}

	err = c.Write([]byte("[]"))
	if err != nil {
		return err //nolint:wrapcheck 	// Error does not need additional wrapping.
	}

	return nil
}

// Get retrieves an item from the cache by ID.
func (c *Cache) Get(identifier string) (*Item, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[identifier]
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrItemNotFound, identifier)
	}

	return item, nil
}

// GetByName retrieves an item from the cache by name.
func (c *Cache) GetByName(name string) (*Item, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, item := range c.items {
		if item.Name == name {
			return item, nil
		}
	}

	return nil, fmt.Errorf("%w: %q", ErrItemNotFound, name)
}

// GetAll retrieves all items from the cache.
func (c *Cache) GetAll() ([]*Item, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	items := make([]*Item, 0, len(c.items))
	for _, item := range c.items {
		items = append(items, item)
	}

	return items, nil
}

// Save stores an item in the cache.
func (c *Cache) Save(item *Item) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[item.ID] = item

	return c.persist()
}

// Delete removes an item from the cache by identifier.
func (c *Cache) Delete(identifier string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.items[identifier]; !ok {
		return fmt.Errorf("%w: %q", ErrItemNotFound, identifier)
	}

	delete(c.items, identifier)

	return c.persist()
}

// Touched returns true if the cache was modified.
func (c *Cache) Touched() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.wasTouched
}

// load reads the cache data from disk.
func (c *Cache) load() error {
	file, err := c.Open()
	if err != nil {
		return err // Error does not need additional wrapping.
	}
	defer file.Close()

	var items []*Item

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&items); err != nil {
		return err // Error does not need additional wrapping.
	}

	c.items = make(map[string]*Item, len(items))
	for _, item := range items {
		c.items[item.ID] = item
	}

	return nil
}

// persist writes the cache data to disk.
func (c *Cache) persist() error {
	file, err := c.OpenForWriting()
	if err != nil {
		return err //nolint:wrapcheck 	// Error does not need additional wrapping.
	}
	defer file.Close()

	items := make([]*Item, 0, len(c.items))
	for _, item := range c.items {
		items = append(items, item)
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	c.wasTouched = true

	return encoder.Encode(items) //nolint:wrapcheck 	// Error does not need additional wrapping.
}

// Set stores an item in the cache (adapter method for Manager interface).
func (c *Cache) Set(id string, item *Item) error {
	item.ID = id
	return c.Save(item)
}

// Remove removes an item from the cache (adapter method for Manager interface).
func (c *Cache) Remove(id string) error {
	return c.Delete(id)
}
