package etcdutil

import (
	"bytes"
	"context"
	"sync"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Kvs struct {
	prefixPath string
	val        map[string][]byte
	mu         sync.RWMutex
}

func NewKvs(ctx context.Context, prefixPath string, client *clientv3.Client) (*Kvs, error) {
	rsp, err := client.Get(ctx, prefixPath, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	kv := &Kvs{val: make(map[string][]byte), prefixPath: prefixPath}
	kv.set(rsp.Kvs)
	return kv, nil
}

func (k *Kvs) ListenPath() string {
	return k.prefixPath
}

func (k *Kvs) Handle(event *clientv3.Event) {
	switch event.Type {
	case mvccpb.PUT:
		k.put(event.Kv)
	case mvccpb.DELETE:
		k.del(event.Kv)
	default:
	}
}

func (k *Kvs) set(data []*mvccpb.KeyValue) {
	if len(data) == 0 {
		return
	}
	k.mu.Lock()
	for _, v := range data {
		i := bytes.LastIndexByte(v.Key, '/')
		k.val[string(v.Key[i+1:])] = v.Value
	}
	k.mu.Unlock()
}

func (k *Kvs) put(data *mvccpb.KeyValue) {
	i := bytes.LastIndexByte(data.Key, '/')
	k.mu.Lock()
	k.val[string(data.Key[i+1:])] = data.Value
	k.mu.Unlock()
}

func (k *Kvs) del(data *mvccpb.KeyValue) {
	i := bytes.LastIndexByte(data.Key, '/')
	k.mu.Lock()
	delete(k.val, BytesToString(data.Key[i+1:]))
	k.mu.Unlock()
}

// Make a copy of the data
func (k *Kvs) Get(key string) ([]byte, bool) {
	k.mu.RLock()
	value, ok := k.val[key]
	k.mu.RUnlock()

	if !ok {
		return nil, false
	}
	r := make([]byte, len(value))
	copy(r, value)
	return r, true
}

type Codec interface {
	Unmarshal(data []byte, v interface{}) error
}

func (k *Kvs) Unmarshal(key string, out any, codec Codec) (bool, error) {
	k.mu.RLock()
	value, ok := k.val[key]
	k.mu.RUnlock()

	if !ok {
		return false, nil
	}

	return true, codec.Unmarshal(value, out)
}
