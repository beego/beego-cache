// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cache

import (
	"context"
	"time"
)

// Cache interface contains all behaviors for cache adapter.
type Cache interface {
	// Get a cached value by key.
	Get(ctx context.Context, key string) (interface{}, error)
	// GetMulti is a batch version of Get.
	GetMulti(ctx context.Context, keys []string) ([]interface{}, error)
	// Put Set a cached value with key and expire time.
	Put(ctx context.Context, key string, val interface{}, timeout time.Duration) error
	// Delete cached value by key.
	// Should not return error if key not found
	Delete(ctx context.Context, key string) error
	// Incr Increment a cached int value by key, as a counter.
	Incr(ctx context.Context, key string) error
	// Decr Decrement a cached int value by key, as a counter.
	Decr(ctx context.Context, key string) error
	// IsExist Check if a cached value exists or not.
	// if key is expired, return (false, nil)
	IsExist(ctx context.Context, key string) (bool, error)
	// ClearAll Clear all cache.
	ClearAll(ctx context.Context) error
}
