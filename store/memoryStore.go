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

func (s *MemoryStore) Get(key string) (Value, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.data[key], nil
}

func (s *MemoryStore) Create(key string, value Value) error {
	s.mutex.Lock()
	s.data[key] = value
	s.mutex.Unlock()
	return nil
}

func (s *MemoryStore) Delete(key string) error {
	s.mutex.Lock()
	delete(s.data, key)
	s.mutex.Unlock()
	return nil
}

func (s *MemoryStore) Increment(key string) error {
	s.mutex.Lock()
	d := s.data[key]
	d.Count++
	s.data[key] = d
	s.mutex.Unlock()
	return nil
}

func (s *MemoryStore) Decrement(key string) error {
	s.mutex.Lock()
	d := s.data[key]
	d.Count--
	s.data[key] = d
	s.mutex.Unlock()
	return nil
}

func (*MemoryStore) Close() error {
	return nil
}
