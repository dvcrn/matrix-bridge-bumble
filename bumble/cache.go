package bumble

import (
	"context"
	"log"
	"sync"
	"time"
)

type Cache struct {
	mu    sync.RWMutex
	items map[string]cacheItem
}

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

func NewCache() *Cache {
	return &Cache{
		items: make(map[string]cacheItem),
	}
}

func (c *Cache) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(1 * time.Minute):
				start := time.Now()
				c.Cleanup()
				log.Printf("cache cleanup took %s", time.Since(start))
			}
		}
	}()
}

func (c *Cache) Set(key string, value interface{}, expiration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	log.Printf("setting cache key %q with expiration of %s", key, expiration)
	c.items[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(expiration),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, ok := c.items[key]
	if !ok {
		log.Printf("cache miss for key %q", key)
		return nil, false
	}
	if time.Now().After(item.expiration) {
		log.Printf("cache hit for key %q but it was expired", key)
		delete(c.items, key)
		return nil, false
	}
	log.Printf("cache hit for key %q", key)
	return item.value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *Cache) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for key, item := range c.items {
		if time.Now().After(item.expiration) {
			delete(c.items, key)
		}
	}
}
