package cache

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIncr(t *testing.T) {
	testCases := []struct {
		name      string
		originVal any
		updateVal any
		afterIncr any
		wantErr   error
	}{
		{
			name:      "int",
			originVal: 1,
			updateVal: 2,
			afterIncr: 1<<(strconv.IntSize-1) - 1,
			wantErr:   ErrIncrementOverflow,
		},
		{
			name:      "int32",
			originVal: int32(1),
			updateVal: int32(2),
			afterIncr: int32(math.MaxInt32),
			wantErr:   ErrIncrementOverflow,
		},
		{
			name:      "int64",
			originVal: int64(1),
			updateVal: int64(2),
			afterIncr: int64(math.MaxInt64),
			wantErr:   ErrIncrementOverflow,
		},
		{
			name:      "uint",
			originVal: uint(1),
			updateVal: uint(2),
			afterIncr: uint(1<<(strconv.IntSize) - 1),
			wantErr:   ErrIncrementOverflow,
		},
		{
			name:      "uint32",
			originVal: uint32(1),
			updateVal: uint32(2),
			afterIncr: uint32(math.MaxUint32),
			wantErr:   ErrIncrementOverflow,
		},
		{
			name:      "uint64",
			originVal: uint64(1),
			updateVal: uint64(2),
			afterIncr: uint64(math.MaxUint64),
			wantErr:   ErrIncrementOverflow,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val, err := incr(tc.originVal)
			assert.Nil(t, err)
			assert.Equal(t, val, tc.updateVal)

			_, err = incr(tc.afterIncr)
			assert.Equal(t, ErrIncrementOverflow, err)
		})
	}
	// other type
	_, err := incr("string")
	assert.Equal(t, ErrNotIntegerType, err)
}

func TestDecr(t *testing.T) {
	testCases := []struct {
		name      string
		originVal any
		updateVal any
		afterDecr any
		wantErr   error
	}{
		{
			name:      "int",
			originVal: 2,
			updateVal: 1,
			afterDecr: -1 << (strconv.IntSize - 1),
			wantErr:   ErrDecrementOverflow,
		},
		{
			name:      "int32",
			originVal: int32(2),
			updateVal: int32(1),
			afterDecr: int32(math.MinInt32),
			wantErr:   ErrDecrementOverflow,
		},
		{
			name:      "int64",
			originVal: int64(2),
			updateVal: int64(1),
			afterDecr: int64(math.MinInt64),
			wantErr:   ErrDecrementOverflow,
		},
		{
			name:      "uint",
			originVal: uint(2),
			updateVal: uint(1),
			afterDecr: uint(0),
			wantErr:   ErrDecrementOverflow,
		},
		{
			name:      "uint32",
			originVal: uint32(2),
			updateVal: uint32(1),
			afterDecr: uint32(0),
			wantErr:   ErrDecrementOverflow,
		},
		{
			name:      "uint64",
			originVal: uint64(2),
			updateVal: uint64(1),
			afterDecr: uint64(0),
			wantErr:   ErrDecrementOverflow,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val, err := decr(tc.originVal)
			assert.Nil(t, err)
			assert.Equal(t, val, tc.updateVal)

			_, err = decr(tc.afterDecr)
			assert.Equal(t, ErrDecrementOverflow, err)
		})
	}
	// other type
	_, err := decr("string")
	assert.Equal(t, ErrNotIntegerType, err)
}
