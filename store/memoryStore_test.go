package store

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMemoryStore(t *testing.T) {
	s := NewMemoryStore()
	assert.NotNil(t, s)
	assert.NotNil(t, s.data)
	assert.NotNil(t, &s.mutex)
}

func TestMemoryStore_Create(t *testing.T) {
	key := "key"
	val := Value{
		Count:     42,
		AccessKey: uuid.New(),
	}

	s := NewMemoryStore()
	s.Create(key, val)

	assert.Equal(t, s.data[key], val)
}

func TestMemoryStore_Get(t *testing.T) {
	key := "key"
	val := Value{
		Count:     42,
		AccessKey: uuid.New(),
	}

	s := NewMemoryStore()
	s.Create(key, val)

	assert.Equal(t, val, s.Get(key))
}

func TestMemoryStore_Delete(t *testing.T) {
	key := "key"
	val := Value{
		Count:     42,
		AccessKey: uuid.New(),
	}

	s := NewMemoryStore()
	s.Create(key, val)
	assert.Equal(t, val, s.Get(key))
	s.Delete(key)
	assert.Equal(t, Value{}, s.Get(key))
}

func TestMemoryStore_Increment(t *testing.T) {
	key := "key"
	val := Value{
		Count:     42,
		AccessKey: uuid.New(),
	}

	s := NewMemoryStore()
	s.Create(key, val)
	s.Increment(key)

	val.Count++

	assert.Equal(t, val, s.Get(key))
}

func TestMemoryStore_Decrement(t *testing.T) {
	key := "key"
	val := Value{
		Count:     42,
		AccessKey: uuid.New(),
	}

	s := NewMemoryStore()
	s.Create(key, val)
	s.Decrement(key)

	val.Count--

	assert.Equal(t, val, s.Get(key))
}
