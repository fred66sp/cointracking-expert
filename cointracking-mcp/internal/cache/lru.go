// Package cache implements the L1 (in-memory LRU+TTL) and L2 (on-disk
// SQLite) cache layers described in SPEC/03-cache-strategy.md.
package cache

import (
	"container/list"
	"sync"
	"time"
)

// Entry is a cached value with its expiry and originating method (the
// method is kept for pattern-based invalidation and stats).
type Entry struct {
	Method    string
	Value     []byte
	ExpiresAt time.Time
}

func (e *Entry) expired(now time.Time) bool {
	return now.After(e.ExpiresAt)
}

type lruItem struct {
	key   string
	entry *Entry
}

// LRU is a fixed-capacity, thread-safe, least-recently-used cache with
// per-entry TTL. Concurrency is guarded by a sync.RWMutex per SPEC 03.
type LRU struct {
	mu      sync.RWMutex
	maxSize int
	items   map[string]*list.Element
	order   *list.List // front = most recently used

	hits   int64
	misses int64
}

func NewLRU(maxSize int) *LRU {
	if maxSize <= 0 {
		maxSize = 1
	}
	return &LRU{
		maxSize: maxSize,
		items:   make(map[string]*list.Element),
		order:   list.New(),
	}
}

// Get returns the entry for key if present and not expired. An expired
// entry is evicted and treated as a miss.
func (c *LRU) Get(key string) (*Entry, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	el, ok := c.items[key]
	if !ok {
		c.misses++
		return nil, false
	}
	item := el.Value.(*lruItem)
	if item.entry.expired(time.Now()) {
		c.removeElementLocked(el)
		c.misses++
		return nil, false
	}
	c.order.MoveToFront(el)
	c.hits++
	return item.entry, true
}

// Set inserts or updates an entry, evicting the least-recently-used item
// if the cache is at capacity.
func (c *LRU) Set(key string, entry *Entry) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if el, ok := c.items[key]; ok {
		el.Value.(*lruItem).entry = entry
		c.order.MoveToFront(el)
		return
	}

	el := c.order.PushFront(&lruItem{key: key, entry: entry})
	c.items[key] = el

	for c.order.Len() > c.maxSize {
		back := c.order.Back()
		if back == nil {
			break
		}
		c.removeElementLocked(back)
	}
}

func (c *LRU) removeElementLocked(el *list.Element) {
	item := el.Value.(*lruItem)
	delete(c.items, item.key)
	c.order.Remove(el)
}

// InvalidateFunc removes all entries whose method or key satisfies match,
// returning the number of entries removed.
func (c *LRU) InvalidateFunc(match func(method, key string) bool) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	removed := 0
	for el := c.order.Front(); el != nil; {
		next := el.Next()
		item := el.Value.(*lruItem)
		if match(item.entry.Method, item.key) {
			c.removeElementLocked(el)
			removed++
		}
		el = next
	}
	return removed
}

// Clear empties the cache and returns the number of entries removed.
func (c *LRU) Clear() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	n := c.order.Len()
	c.items = make(map[string]*list.Element)
	c.order.Init()
	return n
}

func (c *LRU) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.order.Len()
}

type Stats struct {
	Size    int
	MaxSize int
	Hits    int64
	Misses  int64
}

func (c *LRU) Stats() Stats {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return Stats{Size: c.order.Len(), MaxSize: c.maxSize, Hits: c.hits, Misses: c.misses}
}
