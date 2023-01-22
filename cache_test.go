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
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func testMultiTypeIncrDecr(t *testing.T, cache Cache) {
	ctx := context.Background()
	key := "incDecKey"
	testCases := []struct {
		name            string
		beforeIncr      any
		afterIncr       any
		timeoutDuration time.Duration
	}{
		{
			name:            "int",
			beforeIncr:      1,
			afterIncr:       2,
			timeoutDuration: 5 * time.Second,
		},
		{
			name:            "int32",
			beforeIncr:      int32(1),
			afterIncr:       int32(2),
			timeoutDuration: 5 * time.Second,
		},
		{
			name:            "int64",
			beforeIncr:      int64(1),
			afterIncr:       int64(2),
			timeoutDuration: 5 * time.Second,
		},
		{
			name:            "uint",
			beforeIncr:      uint(1),
			afterIncr:       uint(2),
			timeoutDuration: 5 * time.Second,
		},
		{
			name:            "uint32",
			beforeIncr:      uint32(1),
			afterIncr:       uint32(2),
			timeoutDuration: 5 * time.Second,
		},
		{
			name:            "uint64",
			beforeIncr:      uint64(1),
			afterIncr:       uint64(2),
			timeoutDuration: 5 * time.Second,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// testIncrDecr(t, cache, tc.beforeIncr, tc.afterIncr, tc.timeoutDuration)

			assert.Nil(t, cache.Put(ctx, key, tc.beforeIncr, tc.timeoutDuration))
			assert.Nil(t, cache.Incr(ctx, key))

			v, _ := cache.Get(ctx, key)
			assert.Equal(t, tc.afterIncr, v)

			assert.Nil(t, cache.Decr(ctx, key))

			v, _ = cache.Get(ctx, key)
			assert.Equal(t, v, tc.beforeIncr)
			assert.Nil(t, cache.Delete(ctx, key))
		})
	}
}

func testIncrOverFlow(t *testing.T, c Cache, timeout time.Duration) {
	ctx := context.Background()
	key := "incKey"

	assert.Nil(t, c.Put(ctx, key, int64(math.MaxInt64), timeout))
	// int64
	defer func() {
		assert.Nil(t, c.Delete(ctx, key))
	}()
	assert.NotNil(t, c.Incr(ctx, key))
}

func testDecrOverFlow(t *testing.T, c Cache, timeout time.Duration) {
	var err error
	ctx := context.Background()
	key := "decKey"

	// int64
	if err = c.Put(ctx, key, int64(math.MinInt64), timeout); err != nil {
		t.Error("Put Error: ", err.Error())
		return
	}
	defer func() {
		if err = c.Delete(ctx, key); err != nil {
			t.Errorf("Delete error: %s", err.Error())
		}
	}()
	if err = c.Decr(ctx, key); err == nil {
		t.Error("Decr error")
		return
	}
}
