package store

import "errors"

var (
	ErrNoSuchKey = errors.New("no key in store")
)

type Store map[string]interface{}

// Put will overite the key value if the key exists.
func (s Store) Put(key, value string) error {
	s[key] = value

	return nil
}

// Get will get the value of the key if it exists.
func (s Store) Get(key string) (string, error) {
	value, ok := s[key]
	if !ok {
		return "", ErrNoSuchKey
	}

	return value.(string), nil
}

// Delete will delete the key value pair from the store.
func (s Store) Delete(key string) error {
	delete(s, key)

	return nil
}
