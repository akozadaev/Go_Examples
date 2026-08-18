//go:build experiment

package main

import (
	"container/list"
)

// LRUCache представляет собой LRU-кэш
type LRUCache struct {
	capacity int
	cache    map[string]*list.Element
	list     *list.List
}

// Entry представляет элемент кэша (ключ-значение)
type Entry struct {
	key   string
	value interface{}
}

// NewLRUCache создаёт новый LRU-кэш с заданной ёмкостью
func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		list:     list.New(),
	}
}

// Get возвращает значение по ключу и помечает элемент как наиболее используемый
func (c *LRUCache) Get(key string) (interface{}, bool) {
	if elem, ok := c.cache[key]; ok {
		c.list.MoveToFront(elem)
		return elem.Value.(*Entry).value, true
	}
	return nil, false
}

// Add добавляет элемент в кэш
func (c *LRUCache) Add(key string, value interface{}) {
	if elem, ok := c.cache[key]; ok {
		c.list.MoveToFront(elem)
		elem.Value.(*Entry).value = value
		return
	}

	// Добавляем новый элемент
	newElem := c.list.PushFront(&Entry{key: key, value: value})
	c.cache[key] = newElem

	// Удаляем старый элемент, если превышен лимит
	if len(c.cache) > c.capacity {
		c.removeOldest()
	}
}

// Remove удаляет элемент из кэша
func (c *LRUCache) Remove(key string) {
	if elem, ok := c.cache[key]; ok {
		c.list.Remove(elem)
		delete(c.cache, key)
	}
}

// removeOldest удаляет наименее используемый элемент
func (c *LRUCache) removeOldest() {
	elem := c.list.Back()
	if elem != nil {
		c.list.Remove(elem)
		entry := elem.Value.(*Entry)
		delete(c.cache, entry.key)
	}
}

// Size возвращает текущее количество элементов в кэше
func (c *LRUCache) Size() int {
	return len(c.cache)
}
