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
	"errors"
	"testing"
	"time"

	"github.com/beego/beego-cache/v2/redis/internal/mock"
	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestCache_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCmdable := mock.NewMockCmdable(ctrl)

	c := &Cache{
		client: mockCmdable,
		prefix: "testKey",
	}

	ctx := context.Background()

	testCases := []struct {
		name           string
		key            string
		cmdableReturn  interface{}
		expectedResult interface{}
		expectedErr    error
	}{
		{
			name: "Normal case",
			key:  "myKey",
			cmdableReturn: func() any {
				return redis.NewStringResult("myValue", nil)
			}(),
			expectedResult: "myValue",
			expectedErr:    nil,
		},
		{
			name: "Cmdable error case",
			key:  "myKey",
			cmdableReturn: func() any {
				return redis.NewStringResult("", errors.New("some error"))
			}(),
			expectedResult: nil,
			expectedErr:    errors.New("some error"),
		},
	}

	// Iterate through the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			// Set the expectation on mockCmdable
			mockCmdable.EXPECT().
				Get(ctx, c.associate(tc.key)).
				Return(tc.cmdableReturn).
				Times(1)

			result, err := c.Get(ctx, tc.key)

			require.Equal(t, tc.expectedErr, err)
			if err != nil {
				return
			}

			require.Equal(t, tc.expectedResult, result)

		})
	}
}

func TestCache_GetMulti(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCmdable := mock.NewMockCmdable(ctrl)

	c := &Cache{
		client: mockCmdable,
		prefix: "testKey",
	}

	ctx := context.Background()

	testCases := []struct {
		name           string
		keys           []string
		cmdableReturn  interface{}
		expectedResult interface{}
		expectedErr    error
	}{
		{
			name: "Normal case",
			keys: []string{"myKey", "myKey2", "myKey3"},
			cmdableReturn: func() any {
				return redis.NewSliceResult([]interface{}{"myVal", "myVal2", "myVal3"}, nil)
			}(),
			expectedResult: []interface{}{"myVal", "myVal2", "myVal3"},
			expectedErr:    nil,
		},
		{
			name: "Cmdable error case",
			keys: []string{"myKey", "myKey2", "myKey3"},
			cmdableReturn: func() any {
				return redis.NewSliceResult(nil, errors.New("some error"))
			}(),
			expectedResult: nil,
			expectedErr:    errors.New("some error"),
		},
	}

	// Iterate through the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			ks := make([]interface{}, 0, len(tc.keys))
			for _, k := range tc.keys {
				ks = append(ks, c.associate(k))
			}

			// Set the expectation on mockCmdable
			mockCmdable.EXPECT().
				MGet(ctx, ks...).
				Return(tc.cmdableReturn).
				Times(1)

			result, err := c.GetMulti(ctx, tc.keys)

			require.Equal(t, tc.expectedErr, err)
			if err != nil {
				return
			}

			require.Equal(t, tc.expectedResult, result)

		})
	}
}

func TestCache_Put(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCmdable := mock.NewMockCmdable(ctrl)

	c := &Cache{
		client: mockCmdable,
		prefix: "testKey",
	}

	ctx := context.Background()

	testCases := []struct {
		name          string
		key           string
		val           interface{}
		cmdableReturn interface{}
		expectedErr   error
	}{
		{
			name: "Normal case",
			key:  "myKey",
			val:  "myVal",
			cmdableReturn: func() any {
				return redis.NewStatusResult("", nil)
			}(),
			expectedErr: nil,
		},
		{
			name: "Cmdable error case",
			key:  "myKey",
			val:  "myVal",
			cmdableReturn: func() any {
				return redis.NewStatusResult("", errors.New("some error"))
			}(),
			expectedErr: errors.New("some error"),
		},
	}

	// Iterate through the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			// Set the expectation on mockCmdable
			mockCmdable.EXPECT().
				Set(ctx, c.associate(tc.key), tc.val, 10*time.Second).
				Return(tc.cmdableReturn).
				Times(1)

			err := c.Put(ctx, tc.key, tc.val, 10*time.Second)
			require.Equal(t, tc.expectedErr, err)

		})
	}
}

func TestCache_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCmdable := mock.NewMockCmdable(ctrl)

	c := &Cache{
		client: mockCmdable,
		prefix: "testKey",
	}

	ctx := context.Background()

	testCases := []struct {
		name          string
		key           string
		cmdableReturn interface{}
		expectedErr   error
	}{
		{
			name: "Normal case",
			key:  "myKey",
			cmdableReturn: func() any {
				return redis.NewIntResult(1, nil)
			}(),
			expectedErr: nil,
		},
		{
			name: "Cmdable error case",
			key:  "myKey",
			cmdableReturn: func() any {
				return redis.NewIntResult(0, errors.New("some error"))
			}(),
			expectedErr: errors.New("some error"),
		},
	}

	// Iterate through the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			// Set the expectation on mockCmdable
			mockCmdable.EXPECT().
				Del(ctx, c.associate(tc.key)).
				Return(tc.cmdableReturn).
				Times(1)

			err := c.Delete(ctx, tc.key)
			require.Equal(t, tc.expectedErr, err)

		})
	}
}

func TestCache_IsExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCmdable := mock.NewMockCmdable(ctrl)

	c := &Cache{
		client: mockCmdable,
		prefix: "testKey",
	}

	ctx := context.Background()

	testCases := []struct {
		name           string
		key            string
		cmdableReturn  interface{}
		expectedResult bool
		expectedErr    error
	}{
		{
			name: "Normal case",
			key:  "myKey",
			cmdableReturn: func() any {
				return redis.NewIntResult(1, nil)
			}(),
			expectedResult: true,
			expectedErr:    nil,
		},
		{
			name: "Cmdable error case",
			key:  "myKey",
			cmdableReturn: func() any {
				return redis.NewIntResult(0, errors.New("some error"))
			}(),
			expectedResult: false,
			expectedErr:    errors.New("some error"),
		},
	}

	// Iterate through the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			// Set the expectation on mockCmdable
			mockCmdable.EXPECT().
				Exists(ctx, c.associate(tc.key)).
				Return(tc.cmdableReturn).
				Times(1)

			ok, err := c.IsExist(ctx, tc.key)
			require.Equal(t, tc.expectedErr, err)
			require.Equal(t, tc.expectedResult, ok)

		})
	}
}

func TestCache_Incr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCmdable := mock.NewMockCmdable(ctrl)

	c := &Cache{
		client: mockCmdable,
		prefix: "testKey",
	}

	ctx := context.Background()

	testCases := []struct {
		name          string
		key           string
		cmdableReturn interface{}
		expectedErr   error
	}{
		{
			name: "Normal case",
			key:  "myKey",
			cmdableReturn: func() any {
				return redis.NewIntResult(2, nil)
			}(),
			expectedErr: nil,
		},
		{
			name: "Cmdable error case",
			key:  "myKey",
			cmdableReturn: func() any {
				return redis.NewIntResult(0, errors.New("some error"))
			}(),
			expectedErr: errors.New("some error"),
		},
	}

	// Iterate through the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			// Set the expectation on mockCmdable
			mockCmdable.EXPECT().
				Incr(ctx, c.associate(tc.key)).
				Return(tc.cmdableReturn).
				Times(1)

			err := c.Incr(ctx, tc.key)
			require.Equal(t, tc.expectedErr, err)

		})
	}
}

func TestCache_Decr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCmdable := mock.NewMockCmdable(ctrl)

	c := &Cache{
		client: mockCmdable,
		prefix: "testKey",
	}

	ctx := context.Background()

	testCases := []struct {
		name          string
		key           string
		cmdableReturn interface{}
		expectedErr   error
	}{
		{
			name: "Normal case",
			key:  "myKey",
			cmdableReturn: func() any {
				return redis.NewIntResult(1, nil)
			}(),
			expectedErr: nil,
		},
		{
			name: "Cmdable error case",
			key:  "myKey",
			cmdableReturn: func() any {
				return redis.NewIntResult(0, errors.New("some error"))
			}(),
			expectedErr: errors.New("some error"),
		},
	}

	// Iterate through the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			// Set the expectation on mockCmdable
			mockCmdable.EXPECT().
				Decr(ctx, c.associate(tc.key)).
				Return(tc.cmdableReturn).
				Times(1)

			err := c.Decr(ctx, tc.key)
			require.Equal(t, tc.expectedErr, err)

		})
	}
}

func TestCache_Scan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCmdable := mock.NewMockCmdable(ctrl)

	c := &Cache{
		client:    mockCmdable,
		prefix:    "testKey",
		scanCount: 2,
	}

	ctx := context.Background()

	testCases := []struct {
		name           string
		pattern        string
		cursors        []uint64
		cmdableReturn  []interface{}
		expectedResult interface{}
		expectedErr    error
	}{
		{
			name:    "iter 1 time",
			pattern: "testKey:*",
			cursors: []uint64{0},
			cmdableReturn: func() []any {
				return []any{redis.NewScanCmdResult([]string{"testKey:1", "testKey:2"}, 0, nil)}
			}(),
			expectedResult: []string{"testKey:1", "testKey:2"},
			expectedErr:    nil,
		},
		{
			name:    "iter more than 1 time",
			pattern: "testKey:*",
			cursors: []uint64{0, 2},
			cmdableReturn: func() []any {
				return []any{
					redis.NewScanCmdResult([]string{"testKey:1", "testKey:2"}, 2, nil),
					redis.NewScanCmdResult([]string{"testKey:3"}, 0, nil),
				}
			}(),
			expectedResult: []string{"testKey:1", "testKey:2", "testKey:3"},
			expectedErr:    nil,
		},
		{
			name:    "Cmdable error case",
			pattern: "testKey:*",
			cursors: []uint64{0},
			cmdableReturn: func() []any {
				return []any{redis.NewScanCmdResult(nil, 0, errors.New("some error"))}
			}(),
			expectedErr: errors.New("some error"),
		},
	}

	// Iterate through the test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set the expectation on mockCmdable
			for i, cmdableReturn := range tc.cmdableReturn {
				mockCmdable.EXPECT().
					Scan(ctx, tc.cursors[i], tc.pattern, int64(2)).
					Return(cmdableReturn).
					Times(1)
			}

			res, err := c.Scan(ctx, tc.pattern)
			require.Equal(t, tc.expectedErr, err)
			if err != nil {
				return
			}
			require.Equal(t, tc.expectedResult, res)

		})
	}
}
