package cache

import (
	"backend-test-app/internal/entity"
	"sync"
	"time"
)

type cacheEntry struct {
	items     []entity.Item
	expiresAt time.Time
}

type cache struct {
	mu   sync.RWMutex
	data map[string]cacheEntry
	ttl  time.Duration
}

func NewItemsCache(ttl time.Duration) entity.ItemsCache {
	return &cache{
		data: make(map[string]cacheEntry),
		ttl:  ttl,
	}
}

func (c *cache) Get(key string) ([]entity.Item, bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.data[key]
	if !ok || time.Now().After(entry.expiresAt) {
		return nil, false, nil
	}
	return entry.items, true, nil
}

func (c *cache) Set(key string, items []entity.Item) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = cacheEntry{
		items:     items,
		expiresAt: time.Now().Add(c.ttl),
	}
	return nil
}
