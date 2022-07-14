package etcdutil

import (
	"testing"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var client *clientv3.Client

func TestMain(m *testing.M) {
	cli, err := clientv3.New(clientv3.Config{Endpoints: []string{
		"127.0.0.1:2379",
	}})
	if err != nil {
		panic(err)
	}
	client = cli
	m.Run()
	_ = cli.Close()
}
