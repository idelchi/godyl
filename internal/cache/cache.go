// Package cache handles the caching of the downloaded items.
package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/idelchi/godyl/internal/tools/version"
	"github.com/idelchi/godyl/pkg/path/file"
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
	// ID is the unique identifier of the item.
	ID string `json:"id"`

	// Name is the name of the item.
	Name string `json:"name"`

	// Path is the file path of the item.
	Path string `json:"path"`

	// Type is the type of the item.
	Type string `json:"type"`

	// Version of the item.
	Version version.Version `json:"version"`

	// Downloaded is the time when the item was downloaded.
	Downloaded time.Time `json:"downloaded"`

	// Updated is the time when the item was last updated.
	Updated time.Time `json:"updated"`
}

// ErrItemNotFound is returned when an item is not found in the cache.
var ErrItemNotFound = errors.New("item not found")

// Items is a map of items indexed by their names.
type Items map[string]*Item

// AsSlice converts the Items map to a slice of pointers to Item.
func (i Items) AsSlice() []*Item {
	items := make([]*Item, 0, len(i))
	for _, item := range i {
		items = append(items, item)
	}

	return items
}

// Cache is a cache backend that stores data in a JSON file.
type Cache struct {
	file.File // embedded file.File for cache operations

	items      Items        // tracked items in the cache
	mu         sync.RWMutex // mutex for concurrent access
	wasTouched bool         // indicator if the cache was modified
}

// IsEmpty checks if the cache is empty.
func (c *Cache) IsEmpty() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.items) == 0
}

// Load creates or loads the cache file.
func (c *Cache) Load() error {
	// Load existing cache data if the file exists
	if c.Exists() {
		return c.load()
	}

	// Otherwise, create a new cache file
	if err := c.Create(); err != nil {
		return err //nolint:wrapcheck 	// Error does not need additional wrapping.
	}

	// Initialize the cache with an empty array
	if err := c.Write([]byte("[]")); err != nil {
		return err //nolint:wrapcheck 	// Error does not need additional wrapping.
	}

	return nil
}

// Get retrieves items from the cache by ID.
// If no identifiers are provided, it returns all items in the cache.
func (c *Cache) Get(identifiers ...string) ([]*Item, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if len(identifiers) == 0 {
		return c.items.AsSlice(), nil
	}

	items := make([]*Item, 0, len(identifiers))

	var errs []error

	for _, identifier := range identifiers {
		item, err := c.get(identifier)
		if err != nil {
			errs = append(errs, err)
		}

		items = append(items, item)
	}

	return items, errors.Join(errs...)
}

// GetByName retrieves items from the cache by name.
func (c *Cache) GetByName(names ...string) ([]*Item, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	items := make([]*Item, 0, len(names))

	var errs []error

	for _, name := range names {
		item, err := c.getByName(name)
		if err != nil {
			errs = append(errs, err)
		}

		items = append(items, item)
	}

	return items, errors.Join(errs...)
}

// Delete removes items from the cache by identifier.
// If no identifiers are provided, it deletes all items in the cache.
func (c *Cache) Delete(identifiers ...string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(identifiers) == 0 {
		// If no identifiers are provided, delete all items
		c.items = make(Items)

		return c.persist()
	}

	var errs []error

	for _, identifier := range identifiers {
		if err := c.delete(identifier); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

// Add stores items in the cache.
func (c *Cache) Add(items ...*Item) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var errs []error

	for _, item := range items {
		if err := c.add(item); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

// Touched returns true if the cache was modified.
func (c *Cache) Touched() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.wasTouched
}

// DeleteByName removes items from the cache by name.
// If no names are provided, it deletes all items in the cache.
func (c *Cache) DeleteByName(names ...string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(names) == 0 {
		// If no names are provided, delete all items
		c.items = make(Items)

		return c.persist()
	}

	var errs []error

	for _, name := range names {
		if err := c.deleteByName(name); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

// deleteByName removes an item from the cache by name.
func (c *Cache) deleteByName(name string) error {
	for id, item := range c.items {
		if item.Name == name {
			delete(c.items, id)

			return c.persist()
		}
	}

	return fmt.Errorf("%w: %q", ErrItemNotFound, name)
}

// get retrieves an item from the cache by ID.
func (c *Cache) get(identifier string) (*Item, error) {
	item, ok := c.items[identifier]
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrItemNotFound, identifier)
	}

	return item, nil
}

// getByName retrieves an item from the cache by name.
func (c *Cache) getByName(name string) (*Item, error) {
	for _, item := range c.items {
		if item.Name == name {
			return item, nil
		}
	}

	return nil, fmt.Errorf("%w: %q", ErrItemNotFound, name)
}

// Save stores an item in the cache.
func (c *Cache) add(item *Item) error {
	c.items[item.ID] = item

	return c.persist()
}

// delete removes an item from the cache by identifier.
func (c *Cache) delete(identifier string) error {
	if _, ok := c.items[identifier]; !ok {
		return fmt.Errorf("%w: %q", ErrItemNotFound, identifier)
	}

	delete(c.items, identifier)

	return c.persist()
}

// load reads the cache data from disk.
func (c *Cache) load() error {
	file, err := c.Open()
	if err != nil {
		return err //nolint:wrapcheck 	// Error does not need additional wrapping.
	}
	defer file.Close()

	var items []*Item

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&items); err != nil {
		return err //nolint:wrapcheck 	// Error does not need additional wrapping.
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
