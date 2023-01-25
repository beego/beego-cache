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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetString(t *testing.T) {
	testCases := []struct {
		name    string
		value   any
		wantVal any
	}{
		{
			name:    "string",
			value:   "test1",
			wantVal: "test1",
		},
		{
			name:    "bytes",
			value:   []byte("test2"),
			wantVal: "test2",
		},
		{
			name:    "int",
			value:   1,
			wantVal: "1",
		},
		{
			name:    "int64",
			value:   int64(1),
			wantVal: "1",
		},
		{
			name:    "float",
			value:   1.1,
			wantVal: "1.1",
		},
		{
			name:    "nil",
			value:   nil,
			wantVal: "",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.wantVal, GetString(tc.value))
		})
	}
}

func TestGetInt(t *testing.T) {
	testCases := []struct {
		name    string
		value   any
		wantVal any
	}{
		{
			name:    "int",
			value:   1,
			wantVal: 1,
		},
		{
			name:    "int32",
			value:   int32(32),
			wantVal: 32,
		},
		{
			name:    "int64",
			value:   int64(64),
			wantVal: 64,
		},
		{
			name:    "string",
			value:   "128",
			wantVal: 128,
		},
		{
			name:    "nil",
			value:   nil,
			wantVal: 0,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.wantVal, GetInt(tc.value))
		})
	}
}

func TestGetInt64(t *testing.T) {
	testCases := []struct {
		name    string
		value   any
		wantVal any
	}{
		{
			name:    "int",
			value:   1,
			wantVal: int64(1),
		},
		{
			name:    "int32",
			value:   int32(32),
			wantVal: int64(32),
		},
		{
			name:    "int64",
			value:   int64(64),
			wantVal: int64(64),
		},
		{
			name:    "string",
			value:   "128",
			wantVal: int64(128),
		},
		{
			name:    "nil",
			value:   nil,
			wantVal: int64(0),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.wantVal, GetInt64(tc.value))
		})
	}
}

func TestGetFloat64(t *testing.T) {
	testCases := []struct {
		name    string
		value   any
		wantVal any
	}{
		{
			name:    "float32",
			value:   float32(1.11),
			wantVal: 1.11,
		},
		{
			name:    "float",
			value:   1.11,
			wantVal: 1.11,
		},
		{
			name:    "float64",
			value:   1,
			wantVal: float64(1),
		},
		{
			name:    "string",
			value:   "1.11",
			wantVal: 1.11,
		},
		{
			name:    "nil",
			value:   nil,
			wantVal: float64(0),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.wantVal, GetFloat64(tc.value))
		})
	}
}

func TestGetBool(t *testing.T) {
	testCases := []struct {
		name    string
		value   any
		wantVal any
	}{
		{
			name:    "bool",
			value:   true,
			wantVal: true,
		},
		{
			name:    "string",
			value:   "true",
			wantVal: true,
		},
		{
			name:    "nil",
			value:   nil,
			wantVal: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.wantVal, GetBool(tc.value))
		})
	}
}
