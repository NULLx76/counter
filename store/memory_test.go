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
	err := s.Create(key, val)
	assert.NoError(t, err)

	assert.Equal(t, s.data[key], val)
}

func TestMemoryStore_Get(t *testing.T) {
	key := "key"
	val := Value{
		Count:     42,
		AccessKey: uuid.New(),
	}

	s := NewMemoryStore()

	err := s.Create(key, val)
	assert.NoError(t, err)

	gv, err := s.Get(key)
	assert.NoError(t, err)

	assert.Equal(t, val, gv)
}

func TestMemoryStore_Delete(t *testing.T) {
	key := "key"
	val := Value{
		Count:     42,
		AccessKey: uuid.New(),
	}

	s := NewMemoryStore()

	err := s.Create(key, val)
	assert.NoError(t, err)

	nv, err := s.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, val, nv)

	err = s.Delete(key)
	assert.NoError(t, err)

	nv, err = s.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, Value{}, nv)
}

func TestMemoryStore_Increment(t *testing.T) {
	key := "key"
	val := Value{
		Count:     42,
		AccessKey: uuid.New(),
	}

	s := NewMemoryStore()

	err := s.Create(key, val)
	assert.NoError(t, err)

	assert.NoError(t, s.Increment(key))

	val.Count++

	nv, err := s.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, val, nv)
}

func TestMemoryStore_Decrement(t *testing.T) {
	key := "key"
	val := Value{
		Count:     42,
		AccessKey: uuid.New(),
	}

	s := NewMemoryStore()

	err := s.Create(key, val)
	assert.NoError(t, err)

	assert.NoError(t, s.Decrement(key))

	val.Count--

	nv, err := s.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, val, nv)
}
