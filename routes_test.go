package main

import (
	"counter/store"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshal(t *testing.T) {
	key := "/some/key"
	val := store.Value{
		Count:     5,
	}

	marshalled := marshal(key, val)

	var result map[string]int

	err := json.Unmarshal([]byte(marshalled), &result)
	assert.NoError(t, err)

	assert.Equal(t, val.Count, result[key])
}
