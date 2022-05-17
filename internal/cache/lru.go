package cache

import (
	"container/list"
	"sync"
)

type lru struct {
	sync.RWMutex
	elementMap map[string]*list.Element
	list       *list.List
	size, cap  int
}

func NewLRUCache(size int) (*lru, error) {
	em := make(map[string]*list.Element)
	ls := list.New()

	return &lru{
		elementMap: em,
		list:       ls,
		size:       0,
		cap:        size,
	}, nil
}

func (l *lru) Add(key string, value interface{}) bool {
	if l.size < l.cap {

	}

	return false
}
