package cache

import (
	"sync"
	"time"
)

type Cache struct {
	data     map[string]cacheEntry
	mutex    sync.RWMutex
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		data:     make(map[string]cacheEntry),
		interval: interval,
	}
	go cache.reapLoop()
	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	entry, exists := c.data[key]
	if !exists {
		return nil, false
	}
	return entry.val, true
}

// cache management funcs
func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	for range ticker.C {
		c.reap(time.Now().Add(-c.interval))
	}
}

func (c *Cache) reap(threshold time.Time) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for key, entry := range c.data {
		if entry.createdAt.Before(threshold) {
			delete(c.data, key)
		}
	}
}
