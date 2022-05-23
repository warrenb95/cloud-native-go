package cache

import (
	"container/list"
	"errors"
	"sync"

	"github.com/warrenb95/cloud-native-go/internal/model"
)

type lru struct {
	sync.Mutex
	elementMap     map[string]*list.Element
	list           *list.List
	size, capacity int
}

// NewLRUCache will create and return a LRU cache with the provided capacity.
func NewLRUCache(capacity int) (*lru, error) {
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
	}, nil
}

// Put will updated/create the key value to the cache and will return true if the request has replaced an old value.
func (l *lru) Put(value *model.KeyValue) (bool, error) {
	l.Lock()
	defer l.Unlock()

	// check if the key is in the map already
	if elem, ok := l.elementMap[value.Key]; ok {
		// replace the value
		elem.Value = value

		// push to front of list as this was recently used
		l.list.MoveToFront(elem)
		return false, nil
	}

	if l.size < l.capacity {
		elem := l.list.PushFront(value)
		l.elementMap[value.Key] = elem
		l.size++
		return false, nil
	}

	elem := l.list.Back()
	if elem == nil {
		return false, errors.New("least used element is nil")
	}

	kv, ok := elem.Value.(*model.KeyValue)
	if !ok {
		return false, errors.New("element value is an invalid type")
	}

	if _, ok := l.elementMap[kv.Key]; ok {
		delete(l.elementMap, kv.Key)
	}
	l.list.Remove(elem)

	elem.Value = value
	l.list.PushFront(elem)

	l.elementMap[kv.Key] = elem

	return true, nil
}

// Read will get the value from the cache if it exists.
func (l *lru) Read(key string) (interface{}, error) {
	l.Lock()
	defer l.Unlock()

	elem, ok := l.elementMap[key]
	if !ok {
		return nil, errors.New("key value not found")
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
func (l *lru) Delete(key string) {
	l.Lock()
	defer l.Unlock()

	delete(l.elementMap, key)

	l.size--
}
