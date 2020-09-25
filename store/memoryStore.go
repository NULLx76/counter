package store

import "sync"

// MemoryStore is a simple in memory and thread-safe implementation of the Repository interface
type MemoryStore struct {
	data  map[string]Value
	mutex sync.RWMutex
}

// NewMemoryStore creates a new store to save your date
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data:  make(map[string]Value),
		mutex: sync.RWMutex{},
	}
}

func (s *MemoryStore) Get(key string) Value {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.data[key]
}

func (s *MemoryStore) Create(key string, value Value) {
	s.mutex.Lock()
	s.data[key] = value
	s.mutex.Unlock()
}

func (s *MemoryStore) Delete(key string) {
	s.mutex.Lock()
	delete(s.data, key)
	s.mutex.Unlock()
}

func (s *MemoryStore) Increment(key string) {
	s.mutex.Lock()
	d := s.data[key]
	d.Count++
	s.data[key] = d
	s.mutex.Unlock()
}

func (s *MemoryStore) Decrement(key string) {
	s.mutex.Lock()
	d := s.data[key]
	d.Count--
	s.data[key] = d
	s.mutex.Unlock()
}
