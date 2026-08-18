package custom_thread_safe

import (
	"container/list"
	"sync"
)

// Cache — общий интерфейс для кэша
type Cache[K comparable, V any] interface {
	Get(key K) (V, bool)
	Add(key K, value V)
	Remove(key K)
	Size() int
}

// LRUCache — потокобезопасный LRU-кэш
type LRUCache[K comparable, V any] struct {
	capacity int
	cache    map[K]*list.Element
	list     *list.List
	mutex    sync.Mutex
}

// NewLRUCache создаёт новый LRU-кэш с заданной ёмкостью
func NewLRUCache[K comparable, V any](capacity int) *LRUCache[K, V] {
	return &LRUCache[K, V]{
		capacity: capacity,
		cache:    make(map[K]*list.Element),
		list:     list.New(),
	}
}

// Entry — элемент кэша
type entry[K comparable, V any] struct {
	key   K
	value V
}

// Get возвращает значение по ключу
func (c *LRUCache[K, V]) Get(key K) (V, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var zero V
	elem, ok := c.cache[key]
	if !ok {
		return zero, false
	}
	c.list.MoveToFront(elem)
	return elem.Value.(entry[K, V]).value, true
}

// Add добавляет элемент в кэш
func (c *LRUCache[K, V]) Add(key K, value V) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, ok := c.cache[key]; ok {
		c.list.MoveToFront(elem)
		elem.Value = entry[K, V]{key, value}
		return
	}

	newElem := c.list.PushFront(entry[K, V]{key, value})
	c.cache[key] = newElem

	if len(c.cache) > c.capacity {
		c.removeOldest()
	}
}

// Remove удаляет элемент из кэша
func (c *LRUCache[K, V]) Remove(key K) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, ok := c.cache[key]; ok {
		c.list.Remove(elem)
		delete(c.cache, key)
	}
}

// Size возвращает текущее количество элементов
func (c *LRUCache[K, V]) Size() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return len(c.cache)
}

// removeOldest удаляет наименее используемый элемент
func (c *LRUCache[K, V]) removeOldest() {
	elem := c.list.Back()
	if elem != nil {
		c.list.Remove(elem)
		kv := elem.Value.(entry[K, V])
		delete(c.cache, kv.key)
	}
}
