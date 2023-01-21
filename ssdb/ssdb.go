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

package ssdb

import (
	"context"
	"fmt"
	"strings"
	"time"

	cache "github.com/beego/beego-cache/v2"

	berror "github.com/beego/beego-error/v2"
	"github.com/ssdb/gossdb/ssdb"
)

// Cache SSDB adapter
type Cache struct {
	conn     *ssdb.Client
	conninfo []string
}

type CacheOptions func(c *Cache)

// CacheWithConninfo configures conninfo for ssdb
func CacheWithConninfo(conninfo []string) CacheOptions {
	return func(c *Cache) {
		c.conninfo = conninfo
	}
}

// NewSsdbCache creates new ssdb adapter.
func NewSsdbCache(conn *ssdb.Client, opts ...CacheOptions) cache.Cache {
	res := &Cache{
		conn: conn,
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

// Get gets a key's value from memcache.
func (rc *Cache) Get(ctx context.Context, key string) (interface{}, error) {
	value, err := rc.conn.Get(key)
	if err == nil {
		return value, nil
	}
	return nil, berror.Wrapf(err, cache.SsdbCacheCurdFailed, "could not get value, key: %s", key)
}

// GetMulti gets one or keys values from ssdb.
func (rc *Cache) GetMulti(ctx context.Context, keys []string) ([]interface{}, error) {
	size := len(keys)
	values := make([]interface{}, size)

	res, err := rc.conn.Do("multi_get", keys)
	if err != nil {
		return values, berror.Wrapf(err, cache.SsdbCacheCurdFailed, "multi_get failed, key: %v", keys)
	}

	resSize := len(res)
	keyIdx := make(map[string]int)
	for i := 1; i < resSize; i += 2 {
		keyIdx[res[i]] = i
	}

	keysErr := make([]string, 0)
	for i, ki := range keys {
		if _, ok := keyIdx[ki]; !ok {
			keysErr = append(keysErr, fmt.Sprintf("key [%s] error: %s", ki, "key not exist"))
			continue
		}
		values[i] = res[keyIdx[ki]+1]
	}

	if len(keysErr) != 0 {
		return values, berror.Error(cache.MultiGetFailed, strings.Join(keysErr, "; "))
	}

	return values, nil
}

// DelMulti deletes one or more keys from memcache
func (rc *Cache) DelMulti(keys []string) error {
	_, err := rc.conn.Do("multi_del", keys)
	return berror.Wrapf(err, cache.SsdbCacheCurdFailed, "multi_del failed: %v", keys)
}

// Put puts value into memcache.
// value:  must be of type string
func (rc *Cache) Put(ctx context.Context, key string, val interface{}, timeout time.Duration) error {
	v, ok := val.(string)
	if !ok {
		return berror.Errorf(cache.InvalidSsdbCacheValue, "value must be string: %v", val)
	}
	var resp []string
	var err error
	ttl := int(timeout / time.Second)
	if ttl < 0 {
		resp, err = rc.conn.Do("set", key, v)
	} else {
		resp, err = rc.conn.Do("setx", key, v, ttl)
	}
	if err != nil {
		return berror.Wrapf(err, cache.SsdbCacheCurdFailed, "set or setx failed, key: %s", key)
	}
	if len(resp) == 2 && resp[0] == "ok" {
		return nil
	}
	return berror.Errorf(cache.SsdbBadResponse, "the response from SSDB server is invalid: %v", resp)
}

// Delete deletes a value in memcache.
func (rc *Cache) Delete(ctx context.Context, key string) error {
	_, err := rc.conn.Del(key)
	return berror.Wrapf(err, cache.SsdbCacheCurdFailed, "del failed: %s", key)
}

// Incr increases a key's counter.
func (rc *Cache) Incr(ctx context.Context, key string) error {
	_, err := rc.conn.Do("incr", key, 1)
	return berror.Wrapf(err, cache.SsdbCacheCurdFailed, "increase failed: %s", key)
}

// Decr decrements a key's counter.
func (rc *Cache) Decr(ctx context.Context, key string) error {
	_, err := rc.conn.Do("incr", key, -1)
	return berror.Wrapf(err, cache.SsdbCacheCurdFailed, "decrease failed: %s", key)
}

// IsExist checks if a key exists in memcache.
func (rc *Cache) IsExist(ctx context.Context, key string) (bool, error) {
	resp, err := rc.conn.Do("exists", key)
	if err != nil {
		return false, berror.Wrapf(err, cache.SsdbCacheCurdFailed, "exists failed: %s", key)
	}
	if len(resp) == 2 && resp[1] == "1" {
		return true, nil
	}
	return false, nil
}

// ClearAll clears all cached items in ssdb.
// If there are many keys, this method may spent much time.
func (rc *Cache) ClearAll(context.Context) error {
	keyStart, keyEnd, limit := "", "", 50
	resp, err := rc.Scan(keyStart, keyEnd, limit)
	for err == nil {
		size := len(resp)
		if size == 1 {
			return nil
		}
		keys := []string{}
		for i := 1; i < size; i += 2 {
			keys = append(keys, resp[i])
		}
		_, e := rc.conn.Do("multi_del", keys)
		if e != nil {
			return berror.Wrapf(e, cache.SsdbCacheCurdFailed, "multi_del failed: %v", keys)
		}
		keyStart = resp[size-2]
		resp, err = rc.Scan(keyStart, keyEnd, limit)
	}
	return berror.Wrap(err, cache.SsdbCacheCurdFailed, "scan failed")
}

// Scan key all cached in ssdb.
func (rc *Cache) Scan(keyStart string, keyEnd string, limit int) ([]string, error) {
	resp, err := rc.conn.Do("scan", keyStart, keyEnd, limit)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
