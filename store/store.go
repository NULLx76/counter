package store

import "github.com/google/uuid"

//go:generate mockgen -destination mock_store/mock_store.go  . Repository

// Value specifies the structure of each value it contains the key used to modify or delete the key as well
type Value struct {
	// Count is the current counter
	Count int
	// AccessKey is the key required to modify this counter
	AccessKey uuid.UUID
}

// Repository defines the interface for storage backends
type Repository interface {
	// Get gets the value of a specified key
	Get(key string) (Value, error)
	// Create creates the value with a specified key and value
	Create(key string, value Value) error
	// Delete deletes the entry with the specified key
	Delete(key string) error
	// Increment atomically increments the value of the specified key
	Increment(key string) error
	// Decrement atomically decrements the value of the specified key
	Decrement(key string) error
	// Close is the destructor of a repository and should clean up any connection, write back to disk etc.
	Close() error
}
