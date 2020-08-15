package store

import "sync"

type MemoryStore struct {
	data  map[string]int
	mutex sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data:  make(map[string]int),
		mutex: sync.RWMutex{},
	}
}

func (s *MemoryStore) Get(key string) int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.data[key]
}

func (s *MemoryStore) Set(key string, value int) {
	s.mutex.Lock()
	s.data[key] = value
	s.mutex.Unlock()
}

func (s *MemoryStore) Increment(key string) {
	s.mutex.Lock()
	s.data[key] = s.data[key] + 1
	s.mutex.Unlock()
}

func (s *MemoryStore) Decrement(key string) {
	s.mutex.Lock()
	s.data[key] = s.data[key] - 1
	s.mutex.Unlock()
}
