package store

import (
	"sync"

	"github.com/warrenb95/cloud-native-go/internal/model"
)

type Store struct {
	sync.RWMutex
	m map[string]interface{}
}

func New(m map[string]interface{}) *Store {
	return &Store{
		m: m,
	}
}

// Put will overite the key value if the key exists.
func (s *Store) Put(key string, value interface{}) error {
	s.Lock()
	defer s.Unlock()

	s.m[key] = value

	return nil
}

// Get will get the value of the key if it exists.
func (s *Store) Get(key string) (interface{}, error) {
	s.RLock()
	defer s.RUnlock()

	value, ok := s.m[key]
	if !ok {
		return "", model.ErrKeyNotFound
	}

	return value, nil
}

// Delete will delete the key value pair from the store.
func (s *Store) Delete(key string) error {
	s.Lock()
	defer s.Unlock()

	delete(s.m, key)

	return nil
}
