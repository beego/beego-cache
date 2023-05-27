## cache

cache is a Go cache manager. It can use many cache adapters. The repo is inspired by `database/sql` .

## How to install?

	go get github.com/beego/beego/v2/client/cache

## What adapters are supported?

As of now this cache support memory, Memcache and Redis.

## How to use it?

First you must import it

	import (
		"github.com/beego/beego/v2/client/cache"
	)

Then init a Cache (example with memory)

	bm := cache.NewMemoryCache(60)	

Use it like this:

	bm.Put("astaxie", 1, 10 * time.Second)
	bm.Get("astaxie")
	bm.IsExist("astaxie")
	bm.Delete("astaxie")

interval means the gc time. The cache will check at each time interval, whether item has expired.

## Memcache

Memcache use the [gomemcache](http://github.com/bradfitz/gomemcache) client.

## Redis

Redis use the [redigo](http://github.com/gomodule/redigo) client.
