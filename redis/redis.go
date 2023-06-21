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

	"github.com/redis/go-redis/v9"

	cache "github.com/beego/beego-cache/v2"
)

var (
	// DefaultKey defines the collection name of redis for the cache adapter.
	DefaultKey = "beecacheRedis"
)

// Cache is Redis cache adapter.
type Cache struct {
	client    redis.Cmdable // redis client
	key       string
	scanCount int64
}

type CacheOptions func(c *Cache)

// CacheWithScanCount configures scan count for redis
func CacheWithScanCount(count int64) CacheOptions {
	return func(c *Cache) {
		c.scanCount = count
	}
}

// CacheWithKey configures key for redis
func CacheWithKey(key string) CacheOptions {
	return func(c *Cache) {
		c.key = key
	}
}

// NewRedisCache creates a new redis cache with default collection name.
func NewRedisCache(client redis.Cmdable, opts ...CacheOptions) cache.Cache {
	res := &Cache{
		client:    client,
		key:       DefaultKey,
		scanCount: 1024,
	}

	for _, opt := range opts {
		opt(res)
	}
	return res
}

// associate with config key.
func (rc *Cache) associate(originKey interface{}) string {
	return fmt.Sprintf("%s:%s", rc.key, originKey)
}

// Get cache from redis.
func (rc *Cache) Get(ctx context.Context, key string) (interface{}, error) {
	return rc.client.Get(ctx, rc.associate(key)).Result()
}

// GetMulti gets cache from redis.
func (rc *Cache) GetMulti(ctx context.Context, keys []string) ([]interface{}, error) {
	args := make([]string, 0, len(keys))
	for _, key := range keys {
		args = append(args, rc.associate(key))
	}
	return rc.client.MGet(ctx, args...).Result()
}

// Put puts cache into redis.
func (rc *Cache) Put(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	return rc.client.Set(ctx, rc.associate(key), val, timeout).Err()
}

// Delete deletes a key's cache in redis.
func (rc *Cache) Delete(ctx context.Context, key string) error {
	return rc.client.Del(ctx, rc.associate(key)).Err()
}

// IsExist checks cache's existence in redis.
func (rc *Cache) IsExist(ctx context.Context, key string) (bool, error) {
	count, err := rc.client.Exists(ctx, rc.associate(key)).Result()
	return count != 0, err
}

// Incr increases a key's counter in redis.
func (rc *Cache) Incr(ctx context.Context, key string) error {
	return rc.client.Incr(ctx, rc.associate(key)).Err()
}

// Decr decreases a key's counter in redis.
func (rc *Cache) Decr(ctx context.Context, key string) error {
	return rc.client.Decr(ctx, rc.associate(key)).Err()
}

// ClearAll deletes all cache in the redis collection
// Be careful about this method, because it scans all keys and the delete them one by one
func (rc *Cache) ClearAll(ctx context.Context) error {
	cachedKeys, err := rc.Scan(ctx, rc.key+":*")
	if err != nil {
		return err
	}
	if len(cachedKeys) > 0 {
		return rc.client.Del(ctx, cachedKeys...).Err()
	}
	return nil
}

// Scan scans all keys matching a given pattern.
func (rc *Cache) Scan(ctx context.Context, pattern string) ([]string, error) {
	var (
		cursor uint64 = 0 // start
		res    []string
		ks     []string
		err    error
	)
	for {
		ks, cursor, err = rc.client.Scan(ctx, cursor, pattern, rc.scanCount).Result()
		if err != nil {
			return nil, err
		}
		res = append(res, ks...)
		if cursor == 0 { // over
			return res, nil
		}
	}
}
