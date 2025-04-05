package backend

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/idelchi/godyl/internal/cache/cache"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/folder"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// SQLite is a cache backend that stores data in a SQLite database
type SQLite struct {
	db *sql.DB
}

// NewSQLite creates a new SQLite backend
func NewSQLite(path file.File) (*SQLite, error) {
	if err := folder.FromFile(path).Create(); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", path.Path())
	if err != nil {
		return nil, err
	}

	// Create table if it doesn't exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS cache_items (
		name TEXT PRIMARY KEY,
		version TEXT NOT NULL,
		time TIMESTAMP NOT NULL
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		db.Close()
		return nil, err
	}

	return &SQLite{
		db: db,
	}, nil
}

// Get retrieves an item from the cache by name
func (sb *SQLite) Get(name string) (*cache.Item, error) {
	query := "SELECT name, version, time FROM cache_items WHERE name = ?"
	row := sb.db.QueryRow(query, name)

	var item cache.Item
	var timeStr string
	err := row.Scan(&item.Name, &item.Version, &timeStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: %q", cache.ErrItemNotFound, name)
		}
		return nil, err
	}

	// Parse the time string
	item.Time, err = time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

// GetAll retrieves all items from the cache
func (sb *SQLite) GetAll() ([]*cache.Item, error) {
	query := "SELECT name, version, time FROM cache_items"
	rows, err := sb.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*cache.Item
	for rows.Next() {
		var item cache.Item
		var timeStr string
		if err := rows.Scan(&item.Name, &item.Version, &timeStr); err != nil {
			return nil, err
		}

		// Parse the time string
		item.Time, err = time.Parse(time.RFC3339, timeStr)
		if err != nil {
			return nil, err
		}

		items = append(items, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if items == nil {
		return []*cache.Item{}, nil
	}

	return items, nil
}

// Save stores an item in the cache
func (sb *SQLite) Save(item *cache.Item) error {
	query := `INSERT OR REPLACE INTO cache_items (name, version, time) VALUES (?, ?, ?)`
	_, err := sb.db.Exec(query, item.Name, item.Version, item.Time.Format(time.RFC3339))
	return err
}

// Delete removes an item from the cache by name
func (sb *SQLite) Delete(name string) error {
	query := "DELETE FROM cache_items WHERE name = ?"
	result, err := sb.db.Exec(query, name)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("item not found")
	}

	return nil
}

// Close releases any resources held by the backend
func (sb *SQLite) Close() error {
	return sb.db.Close()
}
