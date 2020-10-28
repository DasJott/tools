# CleverReach Cache
A very simple, fast and threadsafe memory cache.

## Environment / crconfig variables
This package uses default values that are automaticly loaded using crconfig.<br>
You can set parameters in environment or in a config file. [See here for more details on crconfig](../crconfig/README.md).<br>

Here are the values and their defaults:
- CACHE_TTL (obj.TTL)<br>
The time the entries may live in seconds.<br>
Default is 600s.<br>

- CACHE_CLEANUP_INTERVAL (obj.CleanupInterval)
The interval time in seconds, the dead entries are actually deleted.<br>
Default is 60s.

## Usage
There are basicly two ways to use this cache.<br>
The easiest way is to use in conjunction with [crconfig](../crconfig/README.md) for getting all settings.<br>
The other way is to define custom config.

### using crconfig

```go
crconfig.Read("myconfig.env")

// Create a cache with config values
mycache := cache.New()

// For new keys, the given function is executed to get the value.
// Within the ttl the function is not executed at all, unless you set the last parameter to true.
// Returned is an interface{} which you can easily cast to expected type.
myval := mycache.Get("mykey", function() interface{} {
    return "this is cached"
}, false).(string)

// Peek just looks for a value by given key and returns wether it could be found
myval, found := mycache.Peek("mykey")

// Set explicitely sets a value to a certain key
mycache.Set("mykey", "new value")

// Delete simply removes the entry from cache
mycache.Delete("mykey")
```

### custom config
```go
// Create a cache with custom values
mycache := cache.Cache{
	// TTL is the default time a cache object is valid, in seconds
	TTL: 600,
	// CleanupInterval is the interval the cache is checked for outdated objects, in seconds. Set 0 to not start interval.
	CleanupInterval: 60,
}
// start interval timer (omit if no cleanup interval is wanted)
mycache.Start()

// For new keys, the given function is executed to get the value.
// Within the ttl the function is not executed at all, unless you set the last parameter to true.
// Returned is an interface{} which you can easily cast to expected type.
myval := mycache.Get("mykey", function() interface{} {
    return "this is cached"
}, false).(string)

// Peek just looks for a value by given key and returns wether it could be found
myval, found := mycache.Peek("mykey")

// Set explicitely sets a value to a certain key
mycache.Set("mykey", "new value")

// Delete simply removes the entry from cache
mycache.Delete("mykey")
```
