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
	"log"
	"os"
	"strings"
	"testing"
	"time"

	cache "github.com/beego/beego-cache/v2"
	"github.com/gomodule/redigo/redis"
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
	t := s.T()
	maxTryCnt := 10

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = s.dsn
	}

	config := fmt.Sprintf(`{"conn": "%s"}`, redisAddr)

	bm, err := cache.NewCache(s.driver, config)

	for err != nil && strings.Contains(cache.InvalidConnection.Desc(), err.Error()) && maxTryCnt > 0 {
		log.Printf("redis 连接异常...")
		time.Sleep(time.Second)

		bm, err = cache.NewCache(s.driver, config)
		maxTryCnt--
	}

	if err != nil {
		t.Fatal(err)
	}
	s.cache = bm

}

type RedisCompositionTestSuite struct {
	Suite
}

func (s *RedisCompositionTestSuite) TearDownTest() {
	// test clear all
	assert.Nil(s.T(), s.cache.ClearAll(context.Background()))
}

func (s *RedisCompositionTestSuite) TestRedisCacheGet() {
	testCases := []struct {
		name            string
		key             string
		value           string
		timeoutDuration time.Duration
		wantErr         error
	}{
		//{
		//	name: "get return err",
		//	key:  "key0",
		//	wantErr: func() error {
		//		err := errors.New("the key not exist")
		//		return berror.Wrapf(err, cache.RedisCacheCurdFailed,
		//			"could not execute this command: %s", "GET")
		//	}(),
		//	timeoutDuration: 1 * time.Second,
		//},
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
			//if err != nil {
			//	assert.EqualError(t, ts.wantErr, err.Error())
			//	return
			//}
			assert.Nil(t, err)
			vs, _ := redis.String(val, err)
			assert.Equal(t, tc.value, vs)
		})
	}
}

func (s *RedisCompositionTestSuite) TestRedisCacheIsExist() {
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

func (s *RedisCompositionTestSuite) TestRedisCacheDelete() {
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
		s.T().Run(tc.name, func(t *testing.T) {
			err := s.cache.Put(context.Background(), tc.key, tc.value, tc.timeoutDuration)
			assert.Nil(t, err)

			err = s.cache.Delete(context.Background(), tc.key)
			assert.Nil(t, err)
		})
	}
}

func (s *RedisCompositionTestSuite) TestRedisCacheGetMulti() {
	testCases := []struct {
		name            string
		keys            []string
		values          []string
		timeoutDuration time.Duration
		wantErr         error
	}{
		//{
		//	name:   "get multi return err",
		//	keys:   []string{"key0", "key1"},
		//	values: []string{"", ""},
		//	wantErr: func() error {
		//		err := errors.New("the key not exist")
		//		return berror.Wrapf(err, cache.RedisCacheCurdFailed,
		//			"could not execute this command: %s", "GET")
		//	}(),
		//	timeoutDuration: 1 * time.Second,
		//},
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
				vs, _ := redis.String(v, err)
				values = append(values, vs)
			}
			assert.Equal(t, tc.values, values)
		})
	}
}

func (s *RedisCompositionTestSuite) TestRedisCacheIncrAndDecr() {
	testCases := []struct {
		name            string
		key             string
		value           int
		timeoutDuration time.Duration
		wantErr         error
	}{
		{
			name:            "incr and decr",
			key:             "key",
			value:           1,
			timeoutDuration: 5 * time.Second,
		},
	}
	for _, tc := range testCases {
		s.T().Run(tc.name, func(t *testing.T) {
			err := s.cache.Put(context.Background(), tc.key, tc.value, tc.timeoutDuration)
			assert.Nil(t, err)

			val, err := s.cache.Get(context.Background(), tc.key)
			assert.Nil(t, err)
			v1, _ := redis.Int(val, err)
			assert.Equal(t, tc.value, v1)

			assert.Nil(t, s.cache.Incr(context.Background(), tc.key))

			val, err = s.cache.Get(context.Background(), tc.key)
			assert.Nil(t, err)
			v2, _ := redis.Int(val, err)
			assert.Equal(t, v1+1, v2)

			assert.Nil(t, s.cache.Decr(context.Background(), tc.key))

			val, err = s.cache.Get(context.Background(), tc.key)
			assert.Nil(t, err)
			v3, _ := redis.Int(val, err)
			assert.Equal(t, tc.value, v3)
		})
	}
}

func (s *RedisCompositionTestSuite) TestCacheScan() {
	t := s.T()
	timeoutDuration := 10 * time.Second

	// insert all
	for i := 0; i < 100; i++ {
		assert.Nil(t, s.cache.Put(context.Background(), fmt.Sprintf("astaxie%d", i), fmt.Sprintf("author%d", i), timeoutDuration))
	}
	time.Sleep(time.Second)
	// scan all for the first time
	keys, err := s.cache.(*Cache).Scan(DefaultKey + ":*")
	assert.Nil(t, err)

	assert.Equal(t, 100, len(keys), "scan all error")

	// clear all
	assert.Nil(t, s.cache.ClearAll(context.Background()))

	// scan all for the second time
	keys, err = s.cache.(*Cache).Scan(DefaultKey + ":*")
	assert.Nil(t, err)
	assert.Equal(t, 0, len(keys))
}

func TestRedisComposition(t *testing.T) {
	suite.Run(t, &RedisCompositionTestSuite{
		Suite{
			driver: "redis",
			dsn:    "127.0.0.1:6379",
		},
	})
}
