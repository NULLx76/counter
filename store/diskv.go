package store

import (
	"encoding/json"
	"github.com/peterbourgon/diskv"
	"os"
	"strings"
)

type DiskvStore struct {
	d *diskv.Diskv
}

func makeKeyPathFriendly(s string) string {
	s = strings.TrimPrefix(s, "/")
	s = strings.ReplaceAll(s, "/", "-")
	return s
}

func NewDiskvStore(path string) *DiskvStore {
	transform := func(s string) []string {
		return strings.Split(s, "-")[0:1]
	}

	d := diskv.New(diskv.Options{
		BasePath:     path,
		Transform:    transform,
		CacheSizeMax: 1024 * 1024,
	})

	return &DiskvStore{d}
}

func (s *DiskvStore) write(key string, value Value) error {
	b, err := json.Marshal(&value)
	if err != nil {
		return err
	}

	return s.d.Write(key, b)
}

func (s *DiskvStore) Create(key string, value Value) error {
	key = makeKeyPathFriendly(key)
	return s.write(key, value)
}

func (s *DiskvStore) Delete(key string) error {
	key = makeKeyPathFriendly(key)
	return s.d.Erase(key)
}

func (s *DiskvStore) Get(key string) (Value, error) {
	key = makeKeyPathFriendly(key)
	val, err := s.d.Read(key)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			return Value{}, nil
		}

		return Value{}, err
	}

	var v Value
	if err := json.Unmarshal(val, &v); err != nil {
		return Value{}, err
	}

	return v, nil
}

func (s *DiskvStore) Increment(key string) error {
	key = makeKeyPathFriendly(key)
	val, err := s.Get(key)
	if err != nil {
		return err
	}

	val.Count++

	return s.write(key, val)
}

func (s *DiskvStore) Decrement(key string) error {
	key = makeKeyPathFriendly(key)
	val, err := s.Get(key)
	if err != nil {
		return err
	}

	val.Count--

	return s.write(key, val)
}

func (s *DiskvStore) Close() error {
	return nil
}
