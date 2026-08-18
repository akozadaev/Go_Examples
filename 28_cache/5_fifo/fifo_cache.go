package main

import "container/list"

// Cache определяет минимальный контракт FIFO-кэша.
type Cache interface {
	Get(key string) (any, bool)
	Add(key string, value any)
	Remove(key string)
	Size() int
}

type entry struct {
	key   string
	value any
}

// FIFOCache — реализация кэша по стратегии "первым пришёл — первым ушёл".
type FIFOCache struct {
	capacity int
	queue    *list.List
	items    map[string]*list.Element
}

// NewFIFOCache создаёт новый FIFO-кэш.
func NewFIFOCache(capacity int) *FIFOCache {
	return &FIFOCache{
		capacity: capacity,
		queue:    list.New(),
		items:    make(map[string]*list.Element),
	}
}

// Add добавляет значение и при переполнении вытесняет самый старый элемент.
func (c *FIFOCache) Add(key string, value any) {
	if c.capacity == 0 {
		return
	}

	if elem, ok := c.items[key]; ok {
		elem.Value.(*entry).value = value
		return
	}

	if c.queue.Len() >= c.capacity {
		c.evict()
	}

	elem := c.queue.PushBack(&entry{key: key, value: value})
	c.items[key] = elem
}

// Get возвращает значение по ключу без изменения порядка.
func (c *FIFOCache) Get(key string) (any, bool) {
	if elem, ok := c.items[key]; ok {
		return elem.Value.(*entry).value, true
	}
	return nil, false
}

// Remove удаляет элемент из кэша.
func (c *FIFOCache) Remove(key string) {
	if elem, ok := c.items[key]; ok {
		c.queue.Remove(elem)
		delete(c.items, key)
	}
}

// Size возвращает текущее количество элементов.
func (c *FIFOCache) Size() int {
	return c.queue.Len()
}

func (c *FIFOCache) evict() {
	elem := c.queue.Front()
	if elem == nil {
		return
	}
	entry := elem.Value.(*entry)
	c.queue.Remove(elem)
	delete(c.items, entry.key)
}
