package store

import (
	"errors"
	"sync"
)

var (
	ErrNoSuchKey = errors.New("no key in store")
)

type TransactionLogger interface {
	WritePut(key string, value interface{})
	WriteDelete(ket string)
}

type Store struct {
	sync.RWMutex
	m  map[string]interface{}
	tr TransactionLogger
}

func New(m map[string]interface{}, tr TransactionLogger) *Store {
	return &Store{
		m:  m,
		tr: tr,
	}
}

// Put will overite the key value if the key exists.
func (s *Store) Put(key, value string) error {
	s.Lock()
	defer s.Unlock()

	s.m[key] = value

	return nil
}

// Get will get the value of the key if it exists.
func (s *Store) Get(key string) (string, error) {
	s.RLock()
	defer s.RUnlock()

	value, ok := s.m[key]
	if !ok {
		return "", ErrNoSuchKey
	}

	return value.(string), nil
}

// Delete will delete the key value pair from the store.
func (s *Store) Delete(key string) error {
	s.Lock()
	defer s.Unlock()

	delete(s.m, key)

	return nil
}
