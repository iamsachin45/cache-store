package main

import (
	"sync"
	"time"
)

type Cache struct {
	data map[string]cacheItem
	mu   sync.RWMutex
	ttl  time.Duration
}

type cacheItem struct {
	value      interface{}
	expiration int64
}

// NewCache creates a cache with TTL
func NewCache(ttl time.Duration) *Cache {
	c := &Cache{
		data: make(map[string]cacheItem),
		ttl:  ttl,
	}
	go c.cleanup() // start background cleanup
	return c
}

// Set a value with expiration
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(c.ttl).UnixNano(),
	}
}

// Get a value (returns nil if expired/not found)
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	item, found := c.data[key]
	c.mu.RUnlock()
	if !found || time.Now().UnixNano() > item.expiration {
		return nil, false
	}
	return item.value, true
}

// Background cleanup for expired items
func (c *Cache) cleanup() {
	for {
		time.Sleep(c.ttl)
		c.mu.Lock()
		for k, v := range c.data {
			if time.Now().UnixNano() > v.expiration {
				delete(c.data, k)
			}
		}
		c.mu.Unlock()
	}
}
