package custom_thread_safe

import (
	"container/list"
	"sync"
)

// Cache определяет контракт для потокобезопасного кэша.
type Cache[K comparable, V any] interface {
	Get(key K) (V, bool)
	Add(key K, value V)
	Remove(key K)
	Size() int
}

type entry[K comparable, V any] struct {
	key   K
	value V
	freq  int
}

// LFUCache реализует потокобезопасный LFU-кэш.
type LFUCache[K comparable, V any] struct {
	capacity int
	minFreq  int

	entries map[K]*list.Element
	freq    map[int]*list.List

	mutex sync.Mutex
}

// NewLFUCache создает новый потокобезопасный LFU-кэш.
func NewLFUCache[K comparable, V any](capacity int) *LFUCache[K, V] {
	return &LFUCache[K, V]{
		capacity: capacity,
		entries:  make(map[K]*list.Element),
		freq:     make(map[int]*list.List),
	}
}

// Get возвращает значение по ключу и обновляет частоту использования.
func (c *LFUCache[K, V]) Get(key K) (V, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var zero V

	elem, ok := c.entries[key]
	if !ok {
		return zero, false
	}

	node := elem.Value.(*entry[K, V])
	c.bump(elem)
	return node.value, true
}

// Add добавляет элемент в кэш, при необходимости вытесняя наименее популярный.
func (c *LFUCache[K, V]) Add(key K, value V) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.capacity == 0 {
		return
	}

	if elem, ok := c.entries[key]; ok {
		node := elem.Value.(*entry[K, V])
		node.value = value
		c.bump(elem)
		return
	}

	if len(c.entries) >= c.capacity {
		c.evict()
	}

	node := &entry[K, V]{
		key:   key,
		value: value,
		freq:  1,
	}

	list := c.ensureList(1)
	elem := list.PushFront(node)
	c.entries[key] = elem
	c.minFreq = 1
}

// Remove удаляет элемент из кэша.
func (c *LFUCache[K, V]) Remove(key K) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	elem, ok := c.entries[key]
	if !ok {
		return
	}

	node := elem.Value.(*entry[K, V])
	list := c.freq[node.freq]
	list.Remove(elem)
	if list.Len() == 0 {
		delete(c.freq, node.freq)
		if c.minFreq == node.freq {
			c.minFreq = c.nextMinFreq()
		}
	}

	delete(c.entries, key)
}

// Size возвращает текущее количество элементов.
func (c *LFUCache[K, V]) Size() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return len(c.entries)
}

func (c *LFUCache[K, V]) bump(elem *list.Element) {
	node := elem.Value.(*entry[K, V])
	currentFreq := node.freq

	currentList := c.freq[currentFreq]
	currentList.Remove(elem)
	if currentList.Len() == 0 {
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

func (c *LFUCache[K, V]) ensureList(freq int) *list.List {
	if l, ok := c.freq[freq]; ok {
		return l
	}
	l := list.New()
	c.freq[freq] = l
	return l
}

func (c *LFUCache[K, V]) evict() {
	list := c.freq[c.minFreq]
	if list == nil {
		return
	}

	elem := list.Back()
	if elem == nil {
		return
	}

	node := elem.Value.(*entry[K, V])
	list.Remove(elem)
	if list.Len() == 0 {
		delete(c.freq, c.minFreq)
	}

	delete(c.entries, node.key)
}

func (c *LFUCache[K, V]) nextMinFreq() int {
	min := 0
	for freq := range c.freq {
		if min == 0 || freq < min {
			min = freq
		}
	}
	return min
}
