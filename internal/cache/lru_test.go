package cache_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/warrenb95/cloud-native-go/internal/cache"
	"github.com/warrenb95/cloud-native-go/internal/model"
)

func Test_lru_Create(t *testing.T) {
	tests := map[string]struct {
		capacity    int
		initValues  []*model.KeyValue
		value       *model.KeyValue
		want        bool
		errContains string
	}{
		"empty cache": {
			capacity: 1,
			value: &model.KeyValue{
				Key:   "key",
				Value: "value",
			},
			want: false,
		},
		"full cache": {
			capacity: 1,
			initValues: []*model.KeyValue{
				{
					Key:   "key1",
					Value: "value1",
				},
			},
			value: &model.KeyValue{
				Key:   "key2",
				Value: "value2",
			},
			want: true,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			lru, err := cache.NewLRUCache(test.capacity)
			require.NoError(t, err)

			for _, val := range test.initValues {
				_, err := lru.Put(val)
				require.NoError(t, err)
			}

			got, err := lru.Put(test.value)
			if test.errContains != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), test.errContains)
			}

			require.NoError(t, err)
			assert.Equal(t, test.want, got)

			expectedSize := math.Min(float64(len(test.initValues)+1), float64(test.capacity))
			assert.Equal(t, expectedSize, float64(lru.Size()))
		})
	}
}

func Test_lru_Get(t *testing.T) {
	tests := map[string]struct {
		capacity    int
		initValues  []*model.KeyValue
		key         string
		want        *model.KeyValue
		errContains string
	}{
		"not found error": {
			capacity: 1,
			key:      "key",
			want: &model.KeyValue{
				Key:   "key",
				Value: "value",
			},
			errContains: "not found",
		},
		"successful": {
			capacity: 1,
			initValues: []*model.KeyValue{
				{
					Key:   "key1",
					Value: "value1",
				},
			},
			key: "key1",
			want: &model.KeyValue{
				Key:   "key1",
				Value: "value1",
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			lru, err := cache.NewLRUCache(test.capacity)
			require.NoError(t, err)

			for _, val := range test.initValues {
				_, err := lru.Put(val)
				require.NoError(t, err)
			}

			got, err := lru.Read(test.key)
			if test.errContains != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), test.errContains)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.want, got)

			expectedSize := math.Min(float64(len(test.initValues)+1), float64(test.capacity))
			assert.Equal(t, expectedSize, float64(lru.Size()))
		})
	}
}

func Test_lru_Delete(t *testing.T) {
	tests := map[string]struct {
		capacity   int
		initValues []*model.KeyValue
		key        string
	}{
		"not found, don't care": {
			capacity: 1,
			key:      "key",
		},
		"successful": {
			capacity: 1,
			initValues: []*model.KeyValue{
				{
					Key:   "key1",
					Value: "value1",
				},
			},
			key: "key1",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			lru, err := cache.NewLRUCache(test.capacity)
			require.NoError(t, err)

			for _, val := range test.initValues {
				_, err := lru.Put(val)
				require.NoError(t, err)
			}

			lru.Delete(test.key)

			_, err = lru.Read(test.key)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "not found")

			expectedSize := math.Min(float64(len(test.initValues)-1), float64(test.capacity))
			assert.Equal(t, expectedSize, float64(lru.Size()))
		})
	}
}
