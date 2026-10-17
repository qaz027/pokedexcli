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
	cachemap map[string]cacheEntry
	// a mutex to protect the map across goroutines
	mu sync.Mutex
}

func NewCache(dur time.Duration) *Cache {
	cacheNew := Cache{cachemap: make(map[string]cacheEntry)}
	go cacheNew.reapLoop(dur)
	return &cacheNew
}

func (c *Cache) Add(key string, val []byte) {
	// lock the mutex
	c.mu.Lock()

	//unlock the mutex
	defer c.mu.Unlock()

	// write the entry into the map
	c.cachemap[key] = cacheEntry{createdAt: time.Now(), val: val}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.cachemap[key]

	if !ok {
		return nil, false
	}

	return entry.val, true

}

func (c *Cache) reapLoop(dur time.Duration) {
	ticker := time.NewTicker(dur)
	for range ticker.C {
		c.mu.Lock()
		for key, entry := range c.cachemap {
			if time.Since(entry.createdAt) > dur {
				delete(c.cachemap, key)
			}
		}
		c.mu.Unlock()
	}
}
