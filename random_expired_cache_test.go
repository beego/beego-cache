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

package cache

import (
	"context"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRandomExpireCache(t *testing.T) {

	bm := NewMemoryCache(20)
	cache := NewRandomExpireCache(bm)
	// should not be nil
	assert.NotNil(t, cache.(*RandomExpireCache).offset)

	timeoutDuration := 3 * time.Second

	if err := cache.Put(context.Background(), "Leon Ding", 22, timeoutDuration); err != nil {
		t.Error("set Error", err)
	}

	// testing random expire cache
	time.Sleep(timeoutDuration + 3 + time.Second)

	if res, _ := cache.IsExist(context.Background(), "Leon Ding"); !res {
		t.Error("check err")
	}

	if v, _ := cache.Get(context.Background(), "Leon Ding"); v.(int) != 22 {
		t.Error("get err")
	}

	assert.Nil(t, cache.Delete(context.Background(), "Leon Ding"))
	res, _ := cache.IsExist(context.Background(), "Leon Ding")
	assert.False(t, res)

	assert.Nil(t, cache.Put(context.Background(), "Leon Ding", "author", timeoutDuration))

	assert.Nil(t, cache.Delete(context.Background(), "astaxie"))
	res, _ = cache.IsExist(context.Background(), "astaxie")
	assert.False(t, res)

	assert.Nil(t, cache.Put(context.Background(), "astaxie", "author", timeoutDuration))

	res, _ = cache.IsExist(context.Background(), "astaxie")
	assert.True(t, res)

	v, _ := cache.Get(context.Background(), "astaxie")
	assert.Equal(t, "author", v)

	assert.Nil(t, cache.Put(context.Background(), "astaxie1", "author1", timeoutDuration))

	res, _ = cache.IsExist(context.Background(), "astaxie1")
	assert.True(t, res)

	vv, _ := cache.GetMulti(context.Background(), []string{"astaxie", "astaxie1"})
	assert.Equal(t, 2, len(vv))
	assert.Equal(t, "author", vv[0])
	assert.Equal(t, "author1", vv[1])

	vv, err := cache.GetMulti(context.Background(), []string{"astaxie0", "astaxie1"})
	assert.Equal(t, 2, len(vv))
	assert.Nil(t, vv[0])
	assert.Equal(t, "author1", vv[1])

	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "key isn't exist"))
}

func TestRandomExpireCacheGet(t *testing.T) {
	testCases := []struct {
		name    string
		key     string
		value   string
		cache   Cache
		wantErr error
	}{
		{
			name:    "key not exist",
			key:     "key0",
			wantErr: ErrKeyNotExist,
			cache: func() Cache {
				bm := NewMemoryCache(20)
				cache := NewRandomExpireCache(bm)
				// should not be nil
				assert.NotNil(t, cache.(*RandomExpireCache).offset)
				return cache
			}(),
		},
		{
			name:  "get val",
			key:   "key2",
			value: "author",
			cache: func() Cache {
				bm := NewMemoryCache(20)
				cache := NewRandomExpireCache(bm)
				// should not be nil
				assert.NotNil(t, cache.(*RandomExpireCache).offset)
				err := cache.Put(context.Background(), "key2", "author", 5*time.Second)
				assert.Nil(t, err)
				return cache
			}(),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val, err := tc.cache.Get(context.Background(), tc.key)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.value, val)
		})
	}
}

func TestRandomExpireCacheIsExist(t *testing.T) {
	bm := NewMemoryCache(20)
	cache := NewRandomExpireCache(bm)
	// should not be nil
	assert.NotNil(t, cache.(*RandomExpireCache).offset)
	testMemoryCacheIsExist(t, cache)
}

func TestRandomExpireCacheDelete(t *testing.T) {
	bm := NewMemoryCache(20)
	cache := NewRandomExpireCache(bm)
	// should not be nil
	assert.NotNil(t, cache.(*RandomExpireCache).offset)
	testMemoryCacheDelete(t, cache)
}

func TestRandomExpireCacheGetMulti(t *testing.T) {
	bm := NewMemoryCache(1)
	cache := NewRandomExpireCache(bm)
	// should not be nil
	assert.NotNil(t, cache.(*RandomExpireCache).offset)
	testMemoryCacheGetMulti(t, cache)
}

func TestWithRandomExpireCacheOffsetFunc(t *testing.T) {
	bm := NewMemoryCache(20)
	magic := -time.Duration(rand.Int())
	cache := NewRandomExpireCache(bm, WithRandomExpireCacheOffsetFunc(func() time.Duration {
		return magic
	}))
	// offset should return the magic value
	assert.Equal(t, magic, cache.(*RandomExpireCache).offset())
}
