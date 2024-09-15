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

// File returns the cache file for the specified cache type.
func File(folder folder.Folder) file.File {
	return folder.WithFile("godyl.json")
}

// New creates a new cache manager with the specified folder as the backend.
func New(folder folder.Folder) *Cache {
	return &Cache{
		File:  File(folder),
		items: make(Items),
	}
}

// Item represents a cache item.
type Item struct {
	// ID is the unique identifier for the item.
	ID string
	// Name is the name of the item.
	Name string
	// Path is the path to the item.
	Path string
	// Time is the time when the item was created or last updated.
	Downloaded time.Time
	// Updated is the time when the item was last updated.
	Updated time.Time
	// Type is the type of the item.
	Type string
	// Commands is the list of commands and patterns to extract the version of the item.
	Version version.Version
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
		return err
	}

	err = c.Create()
	if err != nil {
		return err
	}

	err = c.Write([]byte("[]"))
	if err != nil {
		return err
	}

	return nil
}

// Get retrieves an item from the cache by ID.
func (c *Cache) Get(id string) (*Item, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[id]
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrItemNotFound, id)
	}

	return item, nil
}

// GetByProperty retrieves an item from the cache by a specific property.
func (c *Cache) GetByProperty(property, value string) (*Item, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	accessors := map[string]func(*Item) string{
		"name": func(i *Item) string { return i.Name },
		"path": func(i *Item) string { return i.Path },
	}

	accessor, ok := accessors[property]
	if !ok {
		return nil, fmt.Errorf("invalid property: %s", property)
	}

	for _, item := range c.items {
		if accessor(item) == value {
			return item, nil
		}
	}

	return nil, fmt.Errorf("%w: %q", ErrItemNotFound, value)
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

// Save stores an item in the.
func (c *Cache) Save(item *Item) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[item.ID] = item

	c.wasTouched = true

	return c.persist()
}

// Delete removes an item from the cache by id.
func (c *Cache) Delete(id string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.items[id]; !ok {
		return errors.New("item not found")
	}

	delete(c.items, id)

	c.wasTouched = true

	return c.persist()
}

// load reads the cache data from disk.
func (c *Cache) load() error {
	file, err := c.Open() // Using embedded File.Open() method
	if err != nil {
		return err
	}
	defer file.Close()

	var items []*Item

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&items); err != nil {
		return err
	}

	c.items = make(map[string]*Item, len(items))
	for _, item := range items {
		c.items[item.ID] = item
	}

	return nil
}

// persist writes the cache data to disk.
func (c *Cache) persist() error {
	file, err := c.OpenForWriting() // Using embedded File.OpenForWriting() method
	if err != nil {
		return err
	}
	defer file.Close()

	items := make([]*Item, 0, len(c.items))
	for _, item := range c.items {
		items = append(items, item)
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty formatting

	return encoder.Encode(items)
}

// Touched returns true if the cache was modified.
func (c *Cache) Touched() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.wasTouched
}
