package store

import "github.com/google/uuid"

type Value struct {
	Count     int
	AccessKey uuid.UUID
}

type Repository interface {
	Get(key string) Value
	Set(key string, value Value)
	Increment(key string)
	Decrement(key string)
}
