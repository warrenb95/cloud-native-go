package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore_Put(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := map[string]struct {
		s           Store
		args        args
		expectedErr error
	}{
		"success": {
			s: make(Store),
			args: args{
				key:   "key",
				value: "value",
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.s.Put(test.args.key, test.args.value)
			if test.expectedErr != nil {
				require.EqualError(t, err, test.expectedErr.Error())
				return
			}
			require.NoError(t, err)

			if value, ok := test.s[test.args.key]; !ok {
				t.Fatal("value not found in store")
			} else {
				assert.Equal(t, test.args.value, value)
			}
		})
	}
}

func TestStore_Get(t *testing.T) {
	type args struct {
		key string
	}
	tests := map[string]struct {
		storedValue map[string]string
		s           Store
		args        args
		want        string
		expectedErr error
	}{
		"success": {
			storedValue: map[string]string{
				"test_key": "test_value",
			},
			s: make(Store),
			args: args{
				key: "test_key",
			},
			want: "test_value",
		},
		"no key": {
			s: make(Store),
			args: args{
				key: "test_key",
			},
			expectedErr: ErrNoSuchKey,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			for k, v := range test.storedValue {
				test.s[k] = v
			}

			got, err := test.s.Get(test.args.key)
			if test.expectedErr != nil {
				require.EqualError(t, err, test.expectedErr.Error())
				return
			}
			require.NoError(t, err)

			if got != test.want {
				t.Errorf("Store.Get() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestStore_Delete(t *testing.T) {
	type args struct {
		key string
	}
	tests := map[string]struct {
		storedValues map[string]string
		s            Store
		args         args
		expectedErr  error
	}{
		"success": {
			storedValues: map[string]string{
				"key": "",
			},
			s: make(Store),
			args: args{
				key: "key",
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			for k, v := range test.storedValues {
				test.s[k] = v
			}

			err := test.s.Delete(test.args.key)
			if test.expectedErr != nil {
				require.EqualError(t, err, test.expectedErr.Error())
				return
			}
			require.NoError(t, err)

			if _, ok := test.s[test.args.key]; ok {
				t.Fatalf("key %s should have been deleted", test.args.key)
			}
		})
	}
}
