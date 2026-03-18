package cache_test

import (
	"errors"
	"fmt"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/idelchi/godyl/internal/cache"
	"github.com/idelchi/godyl/internal/tools/version"
	"github.com/idelchi/godyl/pkg/path/file"
)

func newTestCache(t *testing.T) *cache.Cache {
	t.Helper()

	dir := t.TempDir()
	f := file.New(dir, "test-cache.json")
	c := cache.New(f)

	if err := c.Load(); err != nil {
		t.Fatal(err)
	}

	return c
}

func testItem(id, name string) *cache.Item {
	return &cache.Item{
		ID:   id,
		Name: name,
		Path: "/tmp/" + name,
		Type: "github",
	}
}

func TestCacheIsEmpty(t *testing.T) {
	t.Parallel()

	t.Run("fresh cache is empty", func(t *testing.T) {
		t.Parallel()

		c := newTestCache(t)

		if !c.IsEmpty() {
			t.Error("expected fresh cache to be empty")
		}
	})

	t.Run("non-empty after add", func(t *testing.T) {
		t.Parallel()

		c := newTestCache(t)

		if err := c.Add(testItem("id1", "owner/repo")); err != nil {
			t.Fatalf("Add failed: %v", err)
		}

		if c.IsEmpty() {
			t.Error("expected cache to be non-empty after Add")
		}
	})
}

func TestCacheAddAndGet(t *testing.T) {
	t.Parallel()

	t.Run("found", func(t *testing.T) {
		t.Parallel()

		c := newTestCache(t)
		item := testItem("id1", "owner/repo")

		if err := c.Add(item); err != nil {
			t.Fatalf("Add failed: %v", err)
		}

		got, err := c.Get("id1")
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}

		if len(got) != 1 {
			t.Fatalf("expected 1 item, got %d", len(got))
		}

		if got[0].Name != "owner/repo" {
			t.Errorf("expected name %q, got %q", "owner/repo", got[0].Name)
		}
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		c := newTestCache(t)

		_, err := c.Get("nonexistent")
		if !errors.Is(err, cache.ErrItemNotFound) {
			t.Errorf("expected ErrItemNotFound, got %v", err)
		}
	})
}

func TestCacheGetByName(t *testing.T) {
	t.Parallel()

	c := newTestCache(t)

	if err := c.Add(testItem("id1", "owner/repo"), testItem("id2", "other/tool")); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	got, err := c.GetByName("owner/repo")
	if err != nil {
		t.Fatalf("GetByName failed: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 item, got %d", len(got))
	}

	if got[0].ID != "id1" {
		t.Errorf("expected id %q, got %q", "id1", got[0].ID)
	}
}

func TestCacheGetByNameWildcard(t *testing.T) {
	t.Parallel()

	c := newTestCache(t)

	if err := c.Add(testItem("id1", "owner/repo"), testItem("id2", "other/tool")); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	got, err := c.GetByName("owner/*")
	if err != nil {
		t.Fatalf("GetByName with wildcard failed: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 item, got %d", len(got))
	}

	if got[0].Name != "owner/repo" {
		t.Errorf("expected name %q, got %q", "owner/repo", got[0].Name)
	}
}

func TestCacheDelete(t *testing.T) {
	t.Parallel()

	t.Run("deletes existing item", func(t *testing.T) {
		t.Parallel()

		c := newTestCache(t)

		if err := c.Add(testItem("id1", "owner/repo")); err != nil {
			t.Fatalf("Add failed: %v", err)
		}

		if err := c.Delete("id1"); err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		_, err := c.Get("id1")
		if !errors.Is(err, cache.ErrItemNotFound) {
			t.Errorf("expected ErrItemNotFound after Delete, got %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		c := newTestCache(t)

		err := c.Delete("nonexistent")
		if !errors.Is(err, cache.ErrItemNotFound) {
			t.Errorf("expected ErrItemNotFound, got %v", err)
		}
	})

	t.Run("partial delete one exists one does not", func(t *testing.T) {
		t.Parallel()

		c := newTestCache(t)

		if err := c.Add(testItem("id1", "owner/repo")); err != nil {
			t.Fatalf("Add failed: %v", err)
		}

		err := c.Delete("id1", "id2")
		if err == nil {
			t.Fatal("Delete(id1, id2) expected error for missing id2, got nil")
		}

		if !errors.Is(err, cache.ErrItemNotFound) {
			t.Errorf("Delete(id1, id2) expected errors.Is(ErrItemNotFound), got %v", err)
		}

		_, getErr := c.Get("id1")
		if !errors.Is(getErr, cache.ErrItemNotFound) {
			t.Errorf("id1 should have been deleted, but Get returned %v", getErr)
		}
	})
}

func TestCacheDeleteAll(t *testing.T) {
	t.Parallel()

	c := newTestCache(t)

	if err := c.Add(testItem("id1", "tool1"), testItem("id2", "tool2")); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	if err := c.Delete(); err != nil {
		t.Fatalf("Delete (no args) failed: %v", err)
	}

	if !c.IsEmpty() {
		t.Error("expected cache to be empty after Delete with no args")
	}
}

func TestCacheDeleteByName(t *testing.T) {
	t.Parallel()

	c := newTestCache(t)
	item := testItem("id1", "tool1")

	if err := c.Add(item); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	if err := c.DeleteByName("tool1"); err != nil {
		t.Fatalf("DeleteByName failed: %v", err)
	}

	_, err := c.Get("id1")
	if !errors.Is(err, cache.ErrItemNotFound) {
		t.Errorf("expected ErrItemNotFound after DeleteByName, got %v", err)
	}
}

func TestCachePersistence(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	f := file.New(dir, "test-cache.json")

	c1 := cache.New(f)
	if err := c1.Load(); err != nil {
		t.Fatalf("first Load failed: %v", err)
	}

	original := testItem("id1", "owner/repo")

	if err := c1.Add(original); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	c2 := cache.New(f)
	if err := c2.Load(); err != nil {
		t.Fatalf("second Load failed: %v", err)
	}

	got, err := c2.Get("id1")
	if err != nil {
		t.Fatalf("Get on reloaded cache failed: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 item, got %d", len(got))
	}

	if diff := cmp.Diff(original, got[0]); diff != "" {
		t.Errorf("round-trip mismatch (-want +got):\n%s", diff)
	}
}

func TestCacheTouched(t *testing.T) {
	t.Parallel()

	t.Run("fresh cache not touched", func(t *testing.T) {
		t.Parallel()

		c := newTestCache(t)

		if c.Touched() {
			t.Error("expected fresh cache to not be touched")
		}
	})

	t.Run("touched after add", func(t *testing.T) {
		t.Parallel()

		c := newTestCache(t)

		if err := c.Add(testItem("id1", "owner/repo")); err != nil {
			t.Fatalf("Add failed: %v", err)
		}

		if !c.Touched() {
			t.Error("expected cache to be touched after Add")
		}
	})
}

func TestCacheGetAll(t *testing.T) {
	t.Parallel()

	c := newTestCache(t)

	if err := c.Add(testItem("id1", "tool1"), testItem("id2", "tool2")); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	got, err := c.Get()
	if err != nil {
		t.Fatalf("Get (no args) failed: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 items, got %d", len(got))
	}

	ids := make([]string, len(got))
	for i, item := range got {
		ids[i] = item.ID
	}

	slices.Sort(ids)

	if !slices.Equal(ids, []string{"id1", "id2"}) {
		t.Errorf("expected IDs [id1 id2], got %v", ids)
	}
}

func TestCacheDeleteByNameAll(t *testing.T) {
	t.Parallel()

	c := newTestCache(t)

	if err := c.Add(testItem("id1", "tool1"), testItem("id2", "tool2")); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	if err := c.DeleteByName(); err != nil {
		t.Fatalf("DeleteByName (no args) failed: %v", err)
	}

	if !c.IsEmpty() {
		t.Error("expected cache to be empty after DeleteByName with no args")
	}
}

func TestCacheDeleteByNameNotFound(t *testing.T) {
	t.Parallel()

	c := newTestCache(t)

	err := c.DeleteByName("nonexistent")
	if !errors.Is(err, cache.ErrItemNotFound) {
		t.Errorf("expected ErrItemNotFound, got %v", err)
	}
}

func TestCacheAddDuplicateIDOverwrites(t *testing.T) {
	t.Parallel()

	c := newTestCache(t)

	first := testItem("id1", "first-name")
	if err := c.Add(first); err != nil {
		t.Fatalf("first Add failed: %v", err)
	}

	second := &cache.Item{
		ID:   "id1",
		Name: "second-name",
		Path: "/tmp/second-name",
		Type: "gitlab",
	}
	if err := c.Add(second); err != nil {
		t.Fatalf("second Add failed: %v", err)
	}

	got, err := c.Get("id1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 item, got %d", len(got))
	}

	if got[0].Name != "second-name" {
		t.Errorf("expected Name %q after overwrite, got %q", "second-name", got[0].Name)
	}

	if got[0].Type != "gitlab" {
		t.Errorf("expected Type %q after overwrite, got %q", "gitlab", got[0].Type)
	}
}

func TestCacheConcurrentAccess(t *testing.T) {
	t.Parallel()

	c := newTestCache(t)

	if err := c.Add(testItem("id1", "tool1")); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	const goroutines = 10

	var wg sync.WaitGroup

	errs := make(chan error, goroutines*2)

	wg.Add(goroutines * 2)

	for i := range goroutines {
		go func() {
			defer wg.Done()

			if _, err := c.Get("id1"); err != nil {
				errs <- fmt.Errorf("concurrent Get failed: %w", err)
			}
		}()

		go func() {
			defer wg.Done()

			id := fmt.Sprintf("concurrent-%d", i)
			if err := c.Add(testItem(id, "tool-"+id)); err != nil {
				errs <- fmt.Errorf("concurrent Add failed: %w", err)
			}
		}()
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		t.Errorf("%v", err)
	}
}

func TestCacheLoadCorruptJSON(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	f := file.New(dir, "corrupt-cache.json")

	// Write invalid JSON directly to the file so the cache sees it as existing.
	if err := f.Create(); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if err := f.Write([]byte("{this is not valid json")); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	c := cache.New(f)

	if err := c.Load(); err == nil {
		t.Error("expected Load to return an error for corrupt JSON, got nil")
	}
}

func TestCacheAddZeroItems(t *testing.T) {
	t.Parallel()

	c := newTestCache(t)

	// Populate with one item so we have a known starting count.
	if err := c.Add(testItem("id1", "owner/repo")); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	// Call Add with no arguments.
	if err := c.Add(); err != nil {
		t.Errorf("Add() with no args returned unexpected error: %v", err)
	}

	// Length should be unchanged.
	all, err := c.Get()
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if len(all) != 1 {
		t.Errorf("expected cache length 1 after no-op Add, got %d", len(all))
	}
}

func TestCacheGetByNameWildcardNoMatch(t *testing.T) {
	t.Parallel()

	c := newTestCache(t)

	if err := c.Add(testItem("id1", "owner1/tool1"), testItem("id2", "owner2/tool2")); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	_, err := c.GetByName("noowner/*")
	if !errors.Is(err, cache.ErrItemNotFound) {
		t.Errorf("expected ErrItemNotFound for non-matching wildcard, got %v", err)
	}
}

func TestCacheDeleteByNameWildcard(t *testing.T) {
	t.Parallel()

	c := newTestCache(t)

	if err := c.Add(testItem("id1", "owner1/tool1"), testItem("id2", "owner2/tool2")); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	// DeleteByName uses exact name matching (not wildcard); use the exact name to remove owner1/tool1.
	if err := c.DeleteByName("owner1/tool1"); err != nil {
		t.Fatalf("DeleteByName failed: %v", err)
	}

	// owner1/tool1 must be gone.
	_, err := c.Get("id1")
	if !errors.Is(err, cache.ErrItemNotFound) {
		t.Errorf("expected id1 (owner1/tool1) to be deleted, got %v", err)
	}

	// owner2/tool2 must still exist.
	got, err := c.Get("id2")
	if err != nil {
		t.Fatalf("Get id2 failed: %v", err)
	}

	if len(got) != 1 || got[0].Name != "owner2/tool2" {
		t.Errorf("expected owner2/tool2 to remain, got %+v", got)
	}
}

func TestCachePersistenceAllFields(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	f := file.New(dir, "all-fields-cache.json")

	c1 := cache.New(f)
	if err := c1.Load(); err != nil {
		t.Fatalf("first Load failed: %v", err)
	}

	now := time.Now().UTC().Truncate(time.Second)

	original := &cache.Item{
		ID:   "full-id",
		Name: "owner/full-tool",
		Path: "/usr/local/bin/full-tool",
		Type: "gitlab",
		Version: version.Version{
			Version: "v1.2.3",
		},
		Downloaded: now,
		Updated:    now.Add(time.Hour),
	}

	if err := c1.Add(original); err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	// Reload from disk.
	c2 := cache.New(f)
	if err := c2.Load(); err != nil {
		t.Fatalf("second Load failed: %v", err)
	}

	got, err := c2.Get("full-id")
	if err != nil {
		t.Fatalf("Get on reloaded cache failed: %v", err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 item, got %d", len(got))
	}

	restored := got[0]

	if diff := cmp.Diff(original, restored); diff != "" {
		t.Errorf("round-trip mismatch (-want +got):\n%s", diff)
	}
}

func TestCacheConcurrentAccessFinalCount(t *testing.T) {
	t.Parallel()

	c := newTestCache(t)

	// Seed one item before the concurrent phase.
	if err := c.Add(testItem("seed", "seed/tool")); err != nil {
		t.Fatalf("Add seed failed: %v", err)
	}

	const goroutines = 10

	var wg sync.WaitGroup

	errs := make(chan error, goroutines*2)

	wg.Add(goroutines * 2)

	for i := range goroutines {
		go func() {
			defer wg.Done()

			if _, err := c.Get("seed"); err != nil {
				errs <- fmt.Errorf("concurrent Get failed: %w", err)
			}
		}()

		go func() {
			defer wg.Done()

			id := fmt.Sprintf("concurrent-%d", i)
			if err := c.Add(testItem(id, "tool-"+id)); err != nil {
				errs <- fmt.Errorf("concurrent Add failed: %w", err)
			}
		}()
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		t.Errorf("%v", err)
	}

	// After all goroutines complete, the cache must hold the seed item plus
	// one item per goroutine (goroutines distinct IDs).
	all, err := c.Get()
	if err != nil {
		t.Fatalf("final Get failed: %v", err)
	}

	// 1 seed + goroutines concurrent items.
	want := 1 + goroutines
	if len(all) != want {
		t.Errorf("expected %d items after concurrent adds, got %d", want, len(all))
	}
}
