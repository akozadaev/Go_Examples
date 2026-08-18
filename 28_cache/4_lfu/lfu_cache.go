package main

import "container/list"

// Cache описывает минимальный контракт LFU-кэша.
type Cache interface {
	Get(key string) (any, bool)
	Add(key string, value any)
	Remove(key string)
	Size() int
}

type entry struct {
	key   string
	value any
	freq  int
}

// LFUCache — реализация кэша с вытеснением наименее часто используемых элементов.
type LFUCache struct {
	capacity int
	size     int
	minFreq  int
	entries  map[string]*list.Element
	freq     map[int]*list.List
}

// NewLFUCache создает новый LFU-кэш заданной емкости.
func NewLFUCache(capacity int) *LFUCache {
	return &LFUCache{
		capacity: capacity,
		entries:  make(map[string]*list.Element),
		freq:     make(map[int]*list.List),
	}
}

// Get возвращает значение по ключу и обновляет статистику использования.
func (c *LFUCache) Get(key string) (any, bool) {
	if elem, ok := c.entries[key]; ok {
		node := elem.Value.(*entry)
		c.bump(elem)
		return node.value, true
	}
	return nil, false
}

// Add записывает значение в кэш и при необходимости вытесняет наименее популярный элемент.
func (c *LFUCache) Add(key string, value any) {
	if c.capacity == 0 {
		return
	}

	if elem, ok := c.entries[key]; ok {
		node := elem.Value.(*entry)
		node.value = value
		c.bump(elem)
		return
	}

	if c.size >= c.capacity {
		c.evict()
	}

	node := &entry{
		key:   key,
		value: value,
		freq:  1,
	}
	list := c.ensureList(1)
	elem := list.PushFront(node)
	c.entries[key] = elem
	c.minFreq = 1
	c.size++
}

// Remove удаляет элемент из кэша.
func (c *LFUCache) Remove(key string) {
	elem, ok := c.entries[key]
	if !ok {
		return
	}

	node := elem.Value.(*entry)
	freqList := c.freq[node.freq]
	freqList.Remove(elem)
	if freqList.Len() == 0 {
		delete(c.freq, node.freq)
		if c.minFreq == node.freq {
			c.minFreq = c.nextMinFreq()
		}
	}

	delete(c.entries, key)
	c.size--
}

// Size возвращает количество элементов в кэше.
func (c *LFUCache) Size() int {
	return c.size
}

func (c *LFUCache) bump(elem *list.Element) {
	node := elem.Value.(*entry)
	currentFreq := node.freq
	freqList := c.freq[currentFreq]
	freqList.Remove(elem)
	if freqList.Len() == 0 {
		delete(c.freq, currentFreq)
		if c.minFreq == currentFreq {
			c.minFreq++
		}
	}

	node.freq++
	nextList := c.ensureList(node.freq)
	newElem := nextList.PushFront(node)
	c.entries[node.key] = newElem
}

func (c *LFUCache) ensureList(freq int) *list.List {
	if l, ok := c.freq[freq]; ok {
		return l
	}
	l := list.New()
	c.freq[freq] = l
	return l
}

func (c *LFUCache) evict() {
	list := c.freq[c.minFreq]
	if list == nil {
		return
	}

	elem := list.Back()
	if elem == nil {
		return
	}
	node := elem.Value.(*entry)

	list.Remove(elem)
	if list.Len() == 0 {
		delete(c.freq, c.minFreq)
	}

	delete(c.entries, node.key)
	c.size--
}

func (c *LFUCache) nextMinFreq() int {
	min := 0
	for freq := range c.freq {
		if min == 0 || freq < min {
			min = freq
		}
	}
	return min
}
