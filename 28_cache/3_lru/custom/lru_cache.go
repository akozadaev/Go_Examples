package main

import (
	"container/list"
)

// Cache определяет общий контракт для всех кэшей
type Cache interface {
	Get(key string) (any, bool)
	Add(key string, value any)
	Remove(key string)
	Size() int
}

// LRUCache — реализация LRU-кэша
type LRUCache struct {
	capacity int
	cache    map[string]*list.Element
	list     *list.List
}

// Entry — элемент кэша
type Entry struct {
	key   string
	value any
}

// NewLRUCache создаёт новый LRU-кэш
func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		list:     list.New(),
	}
}

// Реализация интерфейса Cache

func (c *LRUCache) Get(key string) (any, bool) {
	if elem, ok := c.cache[key]; ok {
		c.list.MoveToFront(elem)
		return elem.Value.(*Entry).value, true
	}
	return nil, false
}

func (c *LRUCache) Add(key string, value any) {
	if elem, ok := c.cache[key]; ok {
		c.list.MoveToFront(elem)
		elem.Value.(*Entry).value = value
		return
	}

	newElem := c.list.PushFront(&Entry{key: key, value: value})
	c.cache[key] = newElem

	if len(c.cache) > c.capacity {
		c.removeOldest()
	}
}

func (c *LRUCache) Remove(key string) {
	if elem, ok := c.cache[key]; ok {
		c.list.Remove(elem)
		delete(c.cache, key)
	}
}

func (c *LRUCache) Size() int {
	return len(c.cache)
}

// Удаление самого старого элемента
func (c *LRUCache) removeOldest() {
	elem := c.list.Back()
	if elem != nil {
		c.list.Remove(elem)
		entry := elem.Value.(*Entry)
		delete(c.cache, entry.key)
	}
}
