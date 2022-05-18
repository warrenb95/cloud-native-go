package cache_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/warrenb95/cloud-native-go/internal/cache"
	"github.com/warrenb95/cloud-native-go/internal/model"
)

func Test_lru_Add(t *testing.T) {
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
				_, err := lru.Add(val)
				require.NoError(t, err)
			}

			got, err := lru.Add(test.value)
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
