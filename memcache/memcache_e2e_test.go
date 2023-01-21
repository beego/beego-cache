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

package memcache

import (
	"context"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/bradfitz/gomemcache/memcache"

	cache "github.com/beego/beego-cache/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	driver string
	dsn    string
	cache  cache.Cache
}

func (s *Suite) SetupSuite() {
	// t := s.T()
	// maxTryCnt := 10

	//config := fmt.Sprintf(`{"conn": "%s"}`, s.dsn)
	//bm, err := cache.NewCache(s.driver, config)
	//
	//for err != nil && strings.Contains(cache.InvalidConnection.Desc(), err.Error()) && maxTryCnt > 0 {
	//	log.Printf("redis 连接异常...")
	//	time.Sleep(time.Second)
	//
	//	bm, err = cache.NewCache(s.driver, config)
	//	maxTryCnt--
	//}
	//
	//if err != nil {
	//	t.Fatal(err)
	//}

	t := s.T()
	maxTryCnt := 10
	pool := memcache.New(s.dsn)

	// test connection
	err := pool.Ping()
	for err != nil && maxTryCnt > 0 {
		log.Printf("memcache connection exception...")
		err = pool.Ping()
		maxTryCnt--
	}
	if err != nil {
		t.Fatal(err)
	}

	bm := NewMemCacheV2(pool)

	s.cache = bm

}

type MemcacheCompositionTestSuite struct {
	Suite
}

func (s *MemcacheCompositionTestSuite) TearDownTest() {
	// test clear all
	assert.Nil(s.T(), s.cache.ClearAll(context.Background()))
}

func (s *MemcacheCompositionTestSuite) TestMemcacheCacheGet() {
	testCases := []struct {
		name            string
		key             string
		value           string
		timeoutDuration time.Duration
	}{
		{
			name:            "get val",
			key:             "key1",
			value:           "author",
			timeoutDuration: 5 * time.Second,
		},
	}
	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			err := s.cache.Put(context.Background(), tc.key, tc.value, tc.timeoutDuration)
			assert.Nil(t, err)
			time.Sleep(2 * time.Second)

			val, err := s.cache.Get(context.Background(), tc.key)
			assert.Nil(t, err)
			vs := val.([]byte)
			assert.Equal(t, tc.value, string(vs))
		})
	}
}

func (s *MemcacheCompositionTestSuite) TestMemcacheCacheIsExist() {
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
			timeoutDuration: 5 * time.Second,
			isExist:         true,
		},
	}
	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			err := s.cache.Put(context.Background(), tc.key, tc.value, tc.timeoutDuration)
			assert.Nil(t, err)

			time.Sleep(2 * time.Second)

			res, _ := s.cache.IsExist(context.Background(), tc.key)
			assert.Equal(t, res, tc.isExist)
		})
	}
}

func (s *MemcacheCompositionTestSuite) TestMemcacheCacheIncrAndDecr() {
	testCases := []struct {
		name            string
		key             string
		value           string
		timeoutDuration time.Duration
		wantErr         error
	}{
		{
			name:            "incr and decr",
			key:             "key",
			value:           "1",
			timeoutDuration: 5 * time.Second,
		},
	}
	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			err := s.cache.Put(context.Background(), tc.key, tc.value, tc.timeoutDuration)
			assert.Nil(t, err)
			v, _ := strconv.Atoi(tc.value)

			val, err := s.cache.Get(context.Background(), tc.key)
			assert.Nil(t, err)
			v1, _ := strconv.Atoi(string(val.([]byte)))
			assert.Equal(t, v, v1)

			assert.Nil(t, s.cache.Incr(context.Background(), tc.key))

			val, err = s.cache.Get(context.Background(), tc.key)
			assert.Nil(t, err)
			v2, _ := strconv.Atoi(string(val.([]byte)))
			assert.Equal(t, v1+1, v2)

			assert.Nil(t, s.cache.Decr(context.Background(), tc.key))

			val, err = s.cache.Get(context.Background(), tc.key)
			assert.Nil(t, err)
			v3, _ := strconv.Atoi(string(val.([]byte)))
			assert.Equal(t, v, v3)
		})
	}
}

func (s *MemcacheCompositionTestSuite) TestMemcacheCacheDelete() {
	testCases := []struct {
		name            string
		key             string
		value           string
		timeoutDuration time.Duration
		wantErr         error
	}{
		{
			name:            "delete val",
			key:             "key1",
			value:           "author",
			timeoutDuration: 5 * time.Second,
		},
	}
	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			err := s.cache.Put(context.Background(), tc.key, tc.value, tc.timeoutDuration)
			assert.Nil(t, err)
			err = s.cache.Delete(context.Background(), tc.key)
			assert.Nil(t, err)
		})
	}
}

func (s *MemcacheCompositionTestSuite) TestMemcacheCacheGetMulti() {
	testCases := []struct {
		name            string
		keys            []string
		values          []string
		timeoutDuration time.Duration
	}{
		{
			name:            "get multi val",
			keys:            []string{"key2", "key3"},
			values:          []string{"value2", "value3"},
			timeoutDuration: 5 * time.Second,
		},
	}
	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			for idx, key := range tc.keys {
				value := tc.values[idx]
				err := s.cache.Put(context.Background(), key, value, tc.timeoutDuration)
				assert.Nil(t, err)
			}

			time.Sleep(2 * time.Second)

			vals, err := s.cache.GetMulti(context.Background(), tc.keys)
			assert.Nil(t, err)
			values := make([]string, 0, len(tc.values))
			for _, v := range vals {
				vs := v.([]byte)
				values = append(values, string(vs))
			}
			assert.Equal(t, tc.values, values)
		})
	}
}

func TestSsdbComposition(t *testing.T) {
	suite.Run(t, &MemcacheCompositionTestSuite{
		Suite{
			driver: "memcache",
			dsn:    "127.0.0.1:11211",
		},
	})
}
