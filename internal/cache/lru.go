package cache

import (
	"container/list"
	"errors"
	"fmt"
	"sync"

	"github.com/warrenb95/cloud-native-go/internal/model"
)

type Store interface {
	Put(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Delete(key string) error
}

type lru struct {
	sync.Mutex
	elementMap     map[string]*list.Element
	list           *list.List
	size, capacity int

	store Store
}

// NewLRUCache will create and return a LRU cache with the provided capacity.
func NewLRUCache(capacity int, store Store) (*lru, error) {
	if capacity == 0 {
		return nil, errors.New("capacity must be > 0")
	}
	em := make(map[string]*list.Element)
	ls := list.New()

	return &lru{
		elementMap: em,
		list:       ls,
		size:       0,
		capacity:   capacity,
		store:      store,
	}, nil
}

// Put will updated/create the key value to the cache and will return true if the request has replaced an old value.
func (l *lru) Put(key string, value interface{}) error {
	l.Lock()
	defer l.Unlock()

	// check if the key is in the map already
	if elem, ok := l.elementMap[key]; ok {
		// replace the value
		elem.Value = value

		// push to front of list as this was recently used
		l.list.MoveToFront(elem)
		return nil
	}

	err := l.addToCache(&model.KeyValue{
		Key:   key,
		Value: value,
	})
	if err != nil {
		return fmt.Errorf("cannot add to cache: %v", err)
	}

	return l.store.Put(key, value)
}

func (l *lru) addToCache(value *model.KeyValue) error {
	if l.size < l.capacity {
		elem := l.list.PushFront(value)
		l.elementMap[value.Key] = elem
		l.size++
		return nil
	}

	elem := l.list.Back()
	if elem == nil {
		return errors.New("least used element is nil")
	}

	kv, ok := elem.Value.(*model.KeyValue)
	if !ok {
		return errors.New("element value is an invalid type")
	}

	if _, ok := l.elementMap[kv.Key]; ok {
		delete(l.elementMap, kv.Key)
	}
	l.list.Remove(elem)

	elem.Value = value
	l.list.PushFront(elem)

	l.elementMap[kv.Key] = elem

	return nil
}

// Get will get the value from the cache if it exists.
func (l *lru) Get(key string) (interface{}, error) {
	l.Lock()
	defer l.Unlock()

	elem, ok := l.elementMap[key]
	if !ok {
		// add to cache
		val, err := l.store.Get(key)
		if err != nil {
			return nil, err
		}

		err = l.addToCache(&model.KeyValue{
			Key:   key,
			Value: val,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to add key value to cache: %v", err)
		}
	}

	l.list.MoveToFront(elem)
	return elem.Value, nil
}

func (l *lru) Size() int {
	l.Lock()
	defer l.Unlock()
	return l.size
}

// Delete will delete the value if the key exists.
func (l *lru) Delete(key string) error {
	l.Lock()
	defer l.Unlock()

	delete(l.elementMap, key)

	l.size--

	return l.store.Delete(key)
}
