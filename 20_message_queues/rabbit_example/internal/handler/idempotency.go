package handler

import "sync"

type ProcessedStore interface {
	IsProcessed(messageID string) bool
	MarkProcessed(messageID string)
}

// MemoryProcessedStore демонстрирует подавление дубликатов только в течение
// жизни одного процесса. Рабочий сервис должен сохранять запись идемпотентности
// в одной транзакции с изменением бизнес-состояния.
type MemoryProcessedStore struct {
	mu        sync.RWMutex
	processed map[string]struct{}
}

func NewMemoryProcessedStore() *MemoryProcessedStore {
	return &MemoryProcessedStore{processed: make(map[string]struct{})}
}

func (s *MemoryProcessedStore) IsProcessed(messageID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.processed[messageID]
	return ok
}

func (s *MemoryProcessedStore) MarkProcessed(messageID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.processed[messageID] = struct{}{}
}
