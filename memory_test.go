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
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemoryCacheGet(t *testing.T) {
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
				bm := NewMemoryCache(1)
				return bm
			}(),
		},
		{
			name: "key expire",
			key:  "key1",
			cache: func() Cache {
				bm := NewMemoryCache(20)
				err := bm.Put(context.Background(), "key1", "value1", 1*time.Second)
				time.Sleep(2 * time.Second)
				assert.Nil(t, err)
				return bm
			}(),
			wantErr: ErrKeyExpired,
		},
		{
			name:  "get val",
			key:   "key2",
			value: "author",
			cache: func() Cache {
				bm := NewMemoryCache(1)
				err := bm.Put(context.Background(), "key2", "author", 5*time.Second)
				assert.Nil(t, err)
				return bm
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

func TestMemoryCacheIsExist(t *testing.T) {
	cache := NewMemoryCache(1)
	testMemoryCacheIsExist(t, cache)
}

func TestMemoryCacheDelete(t *testing.T) {
	cache := NewMemoryCache(1)
	testMemoryCacheDelete(t, cache)
}

func TestMemoryCacheGetMulti(t *testing.T) {
	cache := NewMemoryCache(1)
	testMemoryCacheGetMulti(t, cache)
}

func TestMemoryCacheIncrAndDecr(t *testing.T) {
	cache := NewMemoryCache(1)
	testMultiTypeIncrDecr(t, cache)
}

func TestMemoryCacheIncrOverFlow(t *testing.T) {
	cache := NewMemoryCache(1)
	testIncrOverFlow(t, cache, time.Second*5)
}

func TestMemoryCacheDecrOverFlow(t *testing.T) {
	cache := NewMemoryCache(1)
	testDecrOverFlow(t, cache, time.Second*5)
}

func TestMemoryCacheConcurrencyIncr(t *testing.T) {
	bm := NewMemoryCache(20)
	err := bm.Put(context.Background(), "edwardhey", 0, time.Second*20)
	assert.Nil(t, err)
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			_ = bm.Incr(context.Background(), "edwardhey")
		}()
	}
	wg.Wait()
	val, _ := bm.Get(context.Background(), "edwardhey")
	if val.(int) != 10 {
		t.Error("Incr err")
	}
}

func testMemoryCacheIsExist(t *testing.T, cache Cache) {
	testCases := []struct {
		name            string
		key             string
		value           string
		timeoutDuration time.Duration
		isExist         bool
	}{
		{
			name:            "not exist",
			key:             "key0",
			value:           "value0",
			timeoutDuration: 1 * time.Second,
		},
		{
			name:            "exist",
			key:             "key1",
			value:           "author",
			timeoutDuration: 12 * time.Second,
			isExist:         true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := cache.Put(context.Background(), tc.key, tc.value, tc.timeoutDuration)
			assert.Nil(t, err)
			time.Sleep(10 * time.Second)
			res, _ := cache.IsExist(context.Background(), tc.key)
			assert.Equal(t, res, tc.isExist)
		})
	}
}

func testMemoryCacheDelete(t *testing.T, cache Cache) {
	testCases := []struct {
		name            string
		key             string
		value           string
		timeoutDuration time.Duration
	}{
		{
			name:            "delete val",
			key:             "key1",
			value:           "author",
			timeoutDuration: 5 * time.Second,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := cache.Put(context.Background(), tc.key, tc.value, tc.timeoutDuration)
			assert.Nil(t, err)
			err = cache.Delete(context.Background(), tc.key)
			assert.Nil(t, err)
		})
	}
}

func testMemoryCacheGetMulti(t *testing.T, cache Cache) {
	testCases := []struct {
		name            string
		keys            []string
		values          []any
		timeoutDuration time.Duration
	}{
		{
			name:            "key not exist",
			keys:            []string{"key0", "key1"},
			values:          []any{"value0", "value1"},
			timeoutDuration: 1 * time.Second,
		},
		{
			name:            "get multi val",
			keys:            []string{"key2", "key3"},
			values:          []any{"value2", "value3"},
			timeoutDuration: 12 * time.Second,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for idx, key := range tc.keys {
				value := tc.values[idx]
				err := cache.Put(context.Background(), key, value, tc.timeoutDuration)
				assert.Nil(t, err)
			}
			time.Sleep(10 * time.Second)
			vals, err := cache.GetMulti(context.Background(), tc.keys)
			if err != nil {
				assert.ErrorContains(t, err, ErrKeyNotExist.Error())
				return
			}
			values := make([]any, 0, len(tc.values))
			for _, val := range vals {
				values = append(values, val)
			}
			assert.Equal(t, tc.values, values)
		})
	}
}
