package cache

import (
	"sync"
	"time"
)

type Cache struct {
	items sync.Map
}

type cacheItem struct {
	value     interface{}
	expiresAt time.Time
}

func (c *Cache) Set(key string, value interface{}, duration time.Duration) {
	c.items.Store(key, cacheItem{
		value:     value,
		expiresAt: time.Now().Add(duration),
	})
}

func (c *Cache) Get(key string) (interface{}, bool) {
	item, ok := c.items.Load(key)
	if !ok {
		return nil, false
	}

	cacheItem := item.(cacheItem)
	if time.Now().After(cacheItem.expiresAt) {
		c.items.Delete(key)
		return nil, false
	}

	return cacheItem.value, true
}

func (c *Cache) Delete(key string) {
	c.items.Delete(key)
}

func (c *Cache) Cleanup() {
	c.items.Range(func(key, value interface{}) bool {
		cacheItem := value.(cacheItem)
		if time.Now().After(cacheItem.expiresAt) {
			c.items.Delete(key)
		}
		return true
	})
}
