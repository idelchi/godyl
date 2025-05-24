// Package cache provides caching functionality for tool installations.
package cache

// Manager defines the interface for cache operations.
type Manager interface {
	// Load loads the cache from storage.
	Load() error
	// Get retrieves an item from the cache by ID.
	Get(id string) (*Item, error)
	// Set stores an item in the cache.
	Set(id string, item *Item) error
	// Remove removes an item from the cache.
	Remove(id string) error
}
