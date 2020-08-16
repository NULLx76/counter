package store

import "github.com/google/uuid"

// Value specifies the structure of each value it contains the key used to modify or delete the key as well
type Value struct {
	// Count is the current counter
	Count     int
	// AccessKey is the key required to modify this counter
	AccessKey uuid.UUID
}

// Repository defines the interface for storage backends
type Repository interface {
	// Get gets the value of a specified key
	Get(key string) Value
	// Set overwrites the value of a specified key with the given value
	Set(key string, value Value)
	// Increment atomically increments the value of the specified key
	Increment(key string)
	// Decrement atomically increments the value of the specified key
	Decrement(key string)
}
