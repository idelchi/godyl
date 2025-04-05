// Package cache provides a caching mechanism with pluggable backends
package cache

import (
	"fmt"
	"time"
)

// Item represents a cache item
type Item struct {
	Name    string    `json:"name"`
	Version string    `json:"version"`
	Time    time.Time `json:"time"`
}

var ErrItemNotFound = fmt.Errorf("item not found")

type Items map[string]*Item

// Backend defines the interface for cache storage backends
type Backend interface {
	// Get retrieves an item from the cache by name
	Get(name string) (*Item, error)

	// GetAll retrieves all items from the cache
	GetAll() ([]*Item, error)

	// Save stores an item in the cache
	Save(item *Item) error

	// Delete removes an item from the cache by name
	Delete(name string) error

	// Close releases any resources held by the backend
	Close() error
}

// Cache represents the cache manager
type Cache struct {
	backend Backend
}

// New creates a new cache with the specified backend
func New(backend Backend) *Cache {
	return &Cache{
		backend: backend,
	}
}

// Get retrieves an item from the cache by name
func (c *Cache) Get(name string) (*Item, error) {
	return c.backend.Get(name)
}

// GetAll retrieves all items from the cache
func (c *Cache) GetAll() ([]*Item, error) {
	return c.backend.GetAll()
}

// Save stores an item in the cache
func (c *Cache) Save(name, version string) error {
	item := &Item{
		Name:    name,
		Version: version,
		Time:    time.Now(),
	}
	return c.backend.Save(item)
}

// Delete removes an item from the cache by name
func (c *Cache) Delete(name string) error {
	return c.backend.Delete(name)
}

// Close releases any resources held by the backend
func (c *Cache) Close() error {
	return c.backend.Close()
}
