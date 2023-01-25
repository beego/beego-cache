beego-cache

The independent cache module from Beego.

---
name: Cache Module
sort: 2
---

# Cache Module

Beego's cache module is used for caching data, inspired by `database/sql`. It supports four cache providers: file, memcache, memory and redis. You can install it by:

	github.com/beego/beego-cache

If you use the `memcache` or `redis` provider, you should first install:

	go get -u github.com/beego/beego-cache/memcache

and then import:

	import _ "github.com/beego/beego-cache/memcache"

## Basic Usage

First step is importing the package:

	import (
		cache "github.com/beego/beego-cache"
	)

Then initialize a global variable object:

- memory

  `interval` stands time, which means the cache will be cleared every 60s:

  	bm := cache.NewMemoryCache(60)

- file

  	bm, err := NewFileCache(
  						FileCacheWithCachePath("cache"),
  						FileCacheWithFileSuffix(".bin"),
  						FileCacheWithDirectoryLevel(2),
  						FileCacheWithEmbedExpiry(120))

- redis

  redis uses [redigo](https://github.com/garyburd/redigo/tree/master/redis)

  	dsn := 127.0.0.1:6379
  	password := "123456"
  	dbNum := 0
  	dialFunc := func() (c redis.Conn, err error) {
  		c, err = redis.Dial("tcp", dsn)
  		if err != nil {
  			return nil, berror.Wrapf(err, cache.DialFailed,
  				"could not dial to remote server: %s ", dsn)
  		}
  	
  		if password != "" {
  			if _, err = c.Do("AUTH", password); err != nil {
  				_ = c.Close()
  				return nil, err
  			}
  		}
  	
  		_, selecterr := c.Do("SELECT", dbNum)
  		if selecterr != nil {
  			_ = c.Close()
  			return nil, selecterr
  		}
  		return
  	}
  	// initialize a new pool
  	pool := &redis.Pool{
  		Dial:        dialFunc,
  		MaxIdle:     3,
  		IdleTimeout: 3 * time.Second,
  	}
  	
  	bm := NewRedisCache(pool)

- memcache

  memcache uses [vitess](http://code.google.com/p/vitess/go/memcache)

  	pool := memcache.New(s.dsn)
  	bm := NewMemCache(pool)

Then we can use `bm` to modify the cache:

	bm.Put("astaxie", 1, 10*time.Second)
	bm.Get("astaxie")
	bm.IsExist("astaxie")
	bm.Delete("astaxie")

## Creating your own provider

The cache module uses the Cache interface, so you can create your own cache provider by implementing this interface and registering it.

```go
type Cache interface {
	// Get a cached value by key.
	Get(ctx context.Context, key string) (interface{}, error)
	// GetMulti is a batch version of Get.
	GetMulti(ctx context.Context, keys []string) ([]interface{}, error)
	// Set a cached value with key and expire time.
	Put(ctx context.Context, key string, val interface{}, timeout time.Duration) error
	// Delete cached value by key.
	Delete(ctx context.Context, key string) error
	// Increment a cached int value by key, as a counter.
	Incr(ctx context.Context, key string) error
	// Decrement a cached int value by key, as a counter.
	Decr(ctx context.Context, key string) error
	// Check if a cached value exists or not.
	IsExist(ctx context.Context, key string) (bool, error)
	// Clear all cache.
	ClearAll(ctx context.Context) error
}
```
