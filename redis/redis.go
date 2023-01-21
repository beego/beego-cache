// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package redis

import (
	"context"
	"fmt"
	"time"

	cache "github.com/beego/beego-cache/v2"

	berror "github.com/beego/beego-error/v2"
	"github.com/gomodule/redigo/redis"
)

// DefaultKey defines the collection name of redis for the cache adapter.
var DefaultKey = "beecacheRedis"

// Cache is Redis cache adapter.
type Cache struct {
	p        *redis.Pool // redis connection pool
	conninfo string
	key      string
}

type CacheOptions func(c *Cache)

// CacheWithConninfo configures conninfo for redis
func CacheWithConninfo(conninfo string) CacheOptions {
	return func(c *Cache) {
		c.conninfo = conninfo
	}
}

// CacheWithKey configures key for redis
func CacheWithKey(key string) CacheOptions {
	return func(c *Cache) {
		c.key = key
	}
}

// NewRedisCache creates a new redis cache with default collection name.
func NewRedisCache(pool *redis.Pool, opts ...CacheOptions) cache.Cache {
	res := &Cache{
		p:   pool,
		key: DefaultKey,
	}

	for _, opt := range opts {
		opt(res)
	}
	return res
}

// Execute the redis commands. args[0] must be the key name
func (rc *Cache) do(commandName string, args ...interface{}) (interface{}, error) {
	args[0] = rc.associate(args[0])
	c := rc.p.Get()
	defer func() {
		_ = c.Close()
	}()

	reply, err := c.Do(commandName, args...)
	if err != nil {
		return nil, berror.Wrapf(err, cache.RedisCacheCurdFailed,
			"could not execute this command: %s", commandName)
	}

	return reply, nil
}

// associate with config key.
func (rc *Cache) associate(originKey interface{}) string {
	return fmt.Sprintf("%s:%s", rc.key, originKey)
}

// Get cache from redis.
func (rc *Cache) Get(ctx context.Context, key string) (interface{}, error) {
	if v, err := rc.do("GET", key); err == nil {
		return v, nil
	} else {
		return nil, err
	}
}

// GetMulti gets cache from redis.
func (rc *Cache) GetMulti(ctx context.Context, keys []string) ([]interface{}, error) {
	c := rc.p.Get()
	defer func() {
		_ = c.Close()
	}()
	var args []interface{}
	for _, key := range keys {
		args = append(args, rc.associate(key))
	}
	return redis.Values(c.Do("MGET", args...))
}

// Put puts cache into redis.
func (rc *Cache) Put(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	_, err := rc.do("SETEX", key, int64(timeout/time.Second), val)
	return err
}

// Delete deletes a key's cache in redis.
func (rc *Cache) Delete(ctx context.Context, key string) error {
	_, err := rc.do("DEL", key)
	return err
}

// IsExist checks cache's existence in redis.
func (rc *Cache) IsExist(ctx context.Context, key string) (bool, error) {
	v, err := redis.Bool(rc.do("EXISTS", key))
	if err != nil {
		return false, err
	}
	return v, nil
}

// Incr increases a key's counter in redis.
func (rc *Cache) Incr(ctx context.Context, key string) error {
	_, err := redis.Bool(rc.do("INCRBY", key, 1))
	return err
}

// Decr decreases a key's counter in redis.
func (rc *Cache) Decr(ctx context.Context, key string) error {
	_, err := redis.Bool(rc.do("INCRBY", key, -1))
	return err
}

// ClearAll deletes all cache in the redis collection
// Be careful about this method, because it scans all keys and the delete them one by one
func (rc *Cache) ClearAll(context.Context) error {
	cachedKeys, err := rc.Scan(rc.key + ":*")
	if err != nil {
		return err
	}
	c := rc.p.Get()
	defer func() {
		_ = c.Close()
	}()
	for _, str := range cachedKeys {
		if _, err = c.Do("DEL", str); err != nil {
			return err
		}
	}
	return err
}

// Scan scans all keys matching a given pattern.
func (rc *Cache) Scan(pattern string) (keys []string, err error) {
	c := rc.p.Get()
	defer func() {
		_ = c.Close()
	}()
	var (
		cursor uint64 = 0 // start
		result []interface{}
		list   []string
	)
	for {
		result, err = redis.Values(c.Do("SCAN", cursor, "MATCH", pattern, "COUNT", 1024))
		if err != nil {
			return
		}
		list, err = redis.Strings(result[1], nil)
		if err != nil {
			return
		}
		keys = append(keys, list...)
		cursor, err = redis.Uint64(result[0], nil)
		if err != nil {
			return
		}
		if cursor == 0 { // over
			return
		}
	}
}
