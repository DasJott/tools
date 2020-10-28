package cache

import (
	"sync"
	"time"

	"cleverreach.com/crtools/crconfig"
)

// Cache represents a cache
// It uses sync.Map, which is a mutexed cache and enhances it by a cleanup routine.
type Cache struct {
	// TTL is the default time a cache object is valid, in seconds
	TTL int64
	// CleanupInterval is the interval the cache is checked for outdated objects, in seconds. Set 0 to not start interval.
	CleanupInterval int64

	mutexCache sync.Map
	running    bool
}

// New returns a pointer to a new Cache object, configured by crconfig/environment
func New() *Cache {
	c := &Cache{
		TTL:             crconfig.GetInt("CACHE_TTL", 600),
		CleanupInterval: crconfig.GetInt("CACHE_CLEANUP_INTERVAL", 60),
		// mutexCache:      sync.Map{},
		running: true,
	}

	c.Start()
	return c
}

// Start can be used if you created an own configured instance of Cache, to start the cleanup interval.
func (c *Cache) Start() {
	if c.CleanupInterval > 0 && !c.running {
		go func() {
			for c.running {
				time.Sleep(time.Duration(c.CleanupInterval) * time.Second)
				now := time.Now().Unix()

				c.mutexCache.Range(func(key, val interface{}) bool {
					if v, ok := val.(cacheObject); ok {
						if v.die <= now {
							c.mutexCache.Delete(key)
						}
					}
					return true
				})
			}
		}()
	}
}

// Stop stops the cleanup interval function
func (c *Cache) Stop() {
	c.running = false
}

type cacheObject struct {
	die  int64
	data interface{}
}

// Peek simply gets the value from cache. No default.
func (c *Cache) Peek(key string) (interface{}, bool) {
	now := time.Now().Unix()

	if val, found := c.mutexCache.Load(key); found {
		if v, ok := val.(cacheObject); ok {
			if v.die > now {
				return v.data, true
			}
		}
	}

	return nil, false
}

// Get gets the value according to given key. Generates and stores it, if new.
// skip returns the given default and replaces the value for given key.
func (c *Cache) Get(key string, f func() interface{}, skip bool) interface{} {
	now := time.Now().Unix()

	if !skip {
		if val, found := c.mutexCache.Load(key); found {
			if v, ok := val.(cacheObject); ok {
				if v.die > now {
					return v.data
				}
			}
		}
	}

	val := f()
	c.mutexCache.Store(key, cacheObject{die: now + c.TTL, data: val})
	return val
}

// Set sets a specific value to a specific key
func (c *Cache) Set(key string, val interface{}) {
	now := time.Now().Unix()
	c.mutexCache.Store(key, cacheObject{die: now + c.TTL, data: val})
}

// Delete explicitely deletes a key from cache.
func (c *Cache) Delete(key string) {
	c.mutexCache.Delete(key)
}

// DeleteAll deletes all keys from cache.
func (c *Cache) DeleteAll() {
	c.mutexCache.Range(func(key, val interface{}) bool {
		c.mutexCache.Delete(key)
		return true
	})
}
