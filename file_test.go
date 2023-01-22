// Copyright 2021 beego
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cache

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFileCacheGet(t *testing.T) {
	testCases := []struct {
		name    string
		key     string
		value   string
		cache   Cache
		wantErr error
	}{
		{
			name:  "get val",
			key:   "key1",
			value: "author",
			cache: func() Cache {
				bm, err := NewFileCache(
					FileCacheWithCachePath("cache"),
					FileCacheWithFileSuffix(".bin"),
					FileCacheWithDirectoryLevel(2),
					FileCacheWithEmbedExpiry(0))
				assert.Nil(t, err)
				err = bm.Put(context.Background(), "key1", "author", 5*time.Second)
				assert.Nil(t, err)
				return bm
			}(),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val, err := tc.cache.Get(context.Background(), tc.key)
			assert.Nil(t, err)
			assert.Equal(t, tc.value, val)
		})
	}
	assert.Nil(t, os.RemoveAll("cache"))
}

func TestFileCacheIsExist(t *testing.T) {
	cache, err := NewFileCache(
		FileCacheWithCachePath("cache"),
		FileCacheWithFileSuffix(".bin"),
		FileCacheWithDirectoryLevel(2),
		FileCacheWithEmbedExpiry(0))
	assert.Nil(t, err)
	testCases := []struct {
		name            string
		key             string
		value           string
		timeoutDuration time.Duration
		isExist         bool
	}{
		{
			name:            "expired",
			key:             "key0",
			value:           "value0",
			timeoutDuration: 1 * time.Second,
			isExist:         true,
		},
		{
			name:            "exist",
			key:             "key1",
			value:           "author",
			timeoutDuration: 5 * time.Second,
			isExist:         true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := cache.Put(context.Background(), tc.key, tc.value, tc.timeoutDuration)
			assert.Nil(t, err)
			time.Sleep(2 * time.Second)
			res, _ := cache.IsExist(context.Background(), tc.key)
			assert.Equal(t, res, tc.isExist)
		})
	}
	assert.Nil(t, os.RemoveAll("cache"))
}

func TestFileCacheDelete(t *testing.T) {
	cache, err := NewFileCache(
		FileCacheWithCachePath("cache"),
		FileCacheWithFileSuffix(".bin"),
		FileCacheWithDirectoryLevel(2),
		FileCacheWithEmbedExpiry(0))
	assert.Nil(t, err)
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
	assert.Nil(t, os.RemoveAll("cache"))
}

func TestFileCacheGetMulti(t *testing.T) {
	cache, err := NewFileCache(
		FileCacheWithCachePath("cache"),
		FileCacheWithFileSuffix(".bin"),
		FileCacheWithDirectoryLevel(2),
		FileCacheWithEmbedExpiry(0))
	assert.Nil(t, err)
	testCases := []struct {
		name            string
		keys            []string
		values          []any
		timeoutDuration time.Duration
	}{
		{
			name:            "key expired",
			keys:            []string{"key0", "key1"},
			values:          []any{"value0", "value1"},
			timeoutDuration: 1 * time.Second,
		},
		{
			name:            "get multi val",
			keys:            []string{"key2", "key3"},
			values:          []any{"value2", "value3"},
			timeoutDuration: 5 * time.Second,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for idx, key := range tc.keys {
				value := tc.values[idx]
				err := cache.Put(context.Background(), key, value, tc.timeoutDuration)
				assert.Nil(t, err)
			}
			time.Sleep(2 * time.Second)
			vals, err := cache.GetMulti(context.Background(), tc.keys)
			if err != nil {
				assert.ErrorContains(t, err, ErrKeyExpired.Error())
				return
			}
			values := make([]any, 0, len(tc.values))
			for _, val := range vals {
				values = append(values, val)
			}
			assert.Equal(t, tc.values, values)
		})
	}
	assert.Nil(t, os.RemoveAll("cache"))
}

func TestFileCacheIncrAndDecr(t *testing.T) {
	cache, err := NewFileCache(
		FileCacheWithCachePath("cache"),
		FileCacheWithFileSuffix(".bin"),
		FileCacheWithDirectoryLevel(2),
		FileCacheWithEmbedExpiry(0))
	assert.Nil(t, err)
	testMultiTypeIncrDecr(t, cache)
	assert.Nil(t, os.RemoveAll("cache"))
}

func TestFileCacheIncrOverFlow(t *testing.T) {
	cache, err := NewFileCache(
		FileCacheWithCachePath("cache"),
		FileCacheWithFileSuffix(".bin"),
		FileCacheWithDirectoryLevel(2),
		FileCacheWithEmbedExpiry(0))
	assert.Nil(t, err)
	testIncrOverFlow(t, cache, time.Second*5)
	assert.Nil(t, os.RemoveAll("cache"))
}

func TestFileCacheDecrOverFlow(t *testing.T) {
	cache, err := NewFileCache(
		FileCacheWithCachePath("cache"),
		FileCacheWithFileSuffix(".bin"),
		FileCacheWithDirectoryLevel(2),
		FileCacheWithEmbedExpiry(0))
	assert.Nil(t, err)
	testDecrOverFlow(t, cache, time.Second*5)
	assert.Nil(t, os.RemoveAll("cache"))
}

func TestFileCacheInit(t *testing.T) {
	fc := &FileCache{}
	FileCacheWithCachePath("////aaa")(fc)
	err := fc.Init()
	assert.NotNil(t, err)
	FileCacheWithCachePath(getTestCacheFilePath())(fc)
	err = fc.Init()
	assert.Nil(t, err)
}

func TestFileGetContents(t *testing.T) {
	_, err := FileGetContents("/bin/aaa")
	assert.NotNil(t, err)
	fn := filepath.Join(os.TempDir(), "fileCache.txt")
	f, err := os.Create(fn)
	assert.Nil(t, err)
	_, err = f.WriteString("text")
	assert.Nil(t, err)
	data, err := FileGetContents(fn)
	assert.Nil(t, err)
	assert.Equal(t, "text", string(data))
}

func TestGobEncodeDecode(t *testing.T) {
	_, err := GobEncode(func() {
		fmt.Print("test func")
	})
	assert.NotNil(t, err)
	data, err := GobEncode(&FileCacheItem{
		Data: "hello",
	})
	assert.Nil(t, err)
	err = GobDecode([]byte("wrong data"), &FileCacheItem{})
	assert.NotNil(t, err)
	dci := &FileCacheItem{}
	err = GobDecode(data, dci)
	assert.Nil(t, err)
	assert.Equal(t, "hello", dci.Data)
}

func getTestCacheFilePath() string {
	return filepath.Join(os.TempDir(), "test", "file.txt")
}
