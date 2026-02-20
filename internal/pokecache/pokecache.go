package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	mu      sync.Mutex
	entries map[string]cacheEntry
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		entries: make(map[string]cacheEntry),
	}
	go cache.reapLoop(interval)
	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = cacheEntry{time.Now().UTC(), val}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if value, ok := c.entries[key]; ok {
		return value.val, true
	} else {
		return nil, false
	}
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.reap(interval)
	}
}

func (c *Cache) reap(interval time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	t := time.Now().UTC()
	for k, v := range c.entries {
		if v.createdAt.Before(t.Add(-interval)) {
			delete(c.entries, k)
		}

	}
}
