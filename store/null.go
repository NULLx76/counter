package store

import "go/types"

// nullStore is a store that doesn't do anything, why? That's a good question.
type nullStore types.Nil

func NewNullStore() *nullStore {
	return &nullStore{}
}
func (nullStore) Get(string) (Value, error) {
	return Value{}, nil
}
func (nullStore) Create(string, Value) error {
	return nil
}
func (nullStore) Delete(string) error {
	return nil
}
func (nullStore) Increment(string) error {
	return nil
}
func (nullStore) Decrement(string) error {
	return nil
}
func (nullStore) Close() error {
	return nil
}
