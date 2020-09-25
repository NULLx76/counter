package store

import (
	"context"
	"encoding/json"
	"go.etcd.io/etcd/clientv3"
	"time"
)

type EtcdStore struct {
	cli *clientv3.Client
	ctx context.Context
}

func NewEtcdStore(endpoints []string) (*EtcdStore, error) {
	return NewEtcdStoreFromConfig(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
}

func NewEtcdStoreFromConfig(cfg clientv3.Config) (*EtcdStore, error) {
	cli, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}

	return &EtcdStore{
		cli,
		context.Background(),
	}, nil
}

func (etcd *EtcdStore) put(key string, value Value) error {
	b, err := json.Marshal(&value)
	if err != nil {
		return err
	}

	_, err = etcd.cli.Put(etcd.ctx, key, string(b))
	return err
}

func (etcd *EtcdStore) Create(key string, value Value) error {
	return etcd.put(key, value)
}

func (etcd *EtcdStore) Delete(key string) error {
	_, err := etcd.cli.Delete(etcd.ctx, key)
	return err
}

func (etcd *EtcdStore) Get(key string) (Value, error) {
	gr, err := etcd.cli.Get(etcd.ctx, key)
	if err != nil {
		return Value{}, err
	}
	if gr.Count < 1 {
		return Value{}, nil
	}

	var v Value
	if err := json.Unmarshal(gr.Kvs[0].Value, &v); err != nil {
		return Value{}, err
	}

	return v, nil
}

func (etcd *EtcdStore) Increment(key string) error {
	val, err := etcd.Get(key)
	if err != nil {
		return err
	}

	val.Count++

	return etcd.put(key, val)
}

func (etcd *EtcdStore) Decrement(key string) error {
	val, err := etcd.Get(key)
	if err != nil {
		return err
	}

	val.Count--

	return etcd.put(key, val)
}

func (etcd *EtcdStore) Close() error {
	return etcd.cli.Close()
}
