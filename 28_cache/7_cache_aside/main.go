package main

import (
	"errors"
	"fmt"
	"sync"
)

var errNotFound = errors.New("not found")

type Repository struct {
	mu    sync.RWMutex
	items map[string]string
	reads int
}

func NewRepository(items map[string]string) *Repository {
	return &Repository{items: items}
}

func (r *Repository) Get(key string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.reads++
	value, ok := r.items[key]
	if !ok {
		return "", errNotFound
	}
	return value, nil
}

func (r *Repository) Update(key, value string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[key] = value
}

type Cache struct {
	mu    sync.RWMutex
	items map[string]string
}

func NewCache() *Cache {
	return &Cache{items: make(map[string]string)}
}

func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.items[key]
	return value, ok
}

func (c *Cache) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = value
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

type Service struct {
	repo  *Repository
	cache *Cache
}

func NewService(repo *Repository, cache *Cache) *Service {
	return &Service{repo: repo, cache: cache}
}

// Get реализует cache-aside: сначала кеш, при miss — источник и заполнение кеша.
func (s *Service) Get(key string) (string, error) {
	if value, ok := s.cache.Get(key); ok {
		fmt.Println("cache hit:", key)
		return value, nil
	}

	fmt.Println("cache miss:", key)
	value, err := s.repo.Get(key)
	if err != nil {
		return "", err
	}
	s.cache.Set(key, value) // ленивое заполнение после первого запроса
	return value, nil
}

// Warm заранее загружает выбранные ключи до пользовательских запросов.
func (s *Service) Warm(keys ...string) error {
	for _, key := range keys {
		value, err := s.repo.Get(key)
		if err != nil {
			return fmt.Errorf("warm %q: %w", key, err)
		}
		s.cache.Set(key, value)
	}
	return nil
}

// Update сначала меняет источник истины, затем инвалидирует устаревшую копию.
func (s *Service) Update(key, value string) {
	s.repo.Update(key, value)
	s.cache.Delete(key)
}

func main() {
	repo := NewRepository(map[string]string{
		"user:1": "Алексей",
		"user:2": "Мария",
	})
	service := NewService(repo, NewCache())

	_ = service.Warm("user:1")
	value, _ := service.Get("user:1") // hit после прогрева
	fmt.Println(value)

	value, _ = service.Get("user:2") // miss и ленивое заполнение
	fmt.Println(value)
	_, _ = service.Get("user:2") // hit

	service.Update("user:2", "Мария (обновлено)") // инвалидация
	value, _ = service.Get("user:2")              // miss, затем новое значение
	fmt.Println(value)
}
