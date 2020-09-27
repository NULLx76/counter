package store

import (
	"encoding/json"
	"github.com/dgraph-io/badger/v2"
)

type BadgerStore struct {
	db *badger.DB
}

func NewBadgerStore(path string) (*BadgerStore, error) {
	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		return nil, err
	}

	return &BadgerStore{db}, nil
}

func (b *BadgerStore) set(key string, value Value) error {
	return b.db.Update(func(txn *badger.Txn) error {
		v, err := json.Marshal(&value)
		if err != nil {
			return err
		}

		return txn.Set([]byte(key), v)
	})
}

func (b *BadgerStore) Create(key string, value Value) error {
	return b.set(key, value)
}

func (b *BadgerStore) Get(key string) (v Value, err error) {
	return v, b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return nil
			} else {
				return err
			}
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &v)
		})
	})
}

func (b *BadgerStore) Delete(key string) error {
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

func (b *BadgerStore) Increment(key string) error {
	v, err := b.Get(key)
	if err != nil {
		return err
	}

	v.Count++

	return b.set(key, v)
}

func (b *BadgerStore) Decrement(key string) error {
	v, err := b.Get(key)
	if err != nil {
		return err
	}

	v.Count--

	return b.set(key, v)
}

func (b *BadgerStore) Close() error {
	return b.db.Close()
}
