package etcdutil

import (
	"context"
	"strings"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type CmdHandle func(command string, args ...string) error

type Cmd struct {
	cmdPrefix string
	resPrefix string
	handler   CmdHandle
	client    *clientv3.Client
}

func NewCmd(prefixPath string, handler CmdHandle, client *clientv3.Client) *Cmd {
	return &Cmd{
		cmdPrefix: prefixPath + "/cmd/",
		resPrefix: prefixPath + "/result/",
		handler:   handler,
		client:    client,
	}
}

func (c *Cmd) Publish(ctx context.Context, command string, args ...string) error {
	var buf strings.Builder
	buf.WriteString(c.cmdPrefix)
	buf.WriteString(command)
	for _, arg := range args {
		buf.WriteString("/")
		buf.WriteString(arg)
	}
	return c.put(ctx, buf.String())
}

func (c *Cmd) ListenPath() string {
	return c.cmdPrefix
}

func (c *Cmd) Handle(event *clientv3.Event) {
	if event.Type != mvccpb.PUT {
		return
	}

	rest := strings.TrimPrefix(BytesToString(event.Kv.Key), c.cmdPrefix)
	arr := strings.Split(rest, "/")
	if len(arr) > 0 {
		if c.handler(arr[0], arr[1:]...) == nil {
			ctx, cancel := context.WithTimeout(c.client.Ctx(), time.Second)
			_ = c.put(ctx, c.resPrefix+rest)
			cancel()
		}
	}
}

func (c *Cmd) put(ctx context.Context, str string) error {
	leaseRsp, err := c.client.Grant(ctx, 180)
	if err != nil {
		return err
	}

	_, err = c.client.Put(ctx, str, "", clientv3.WithLease(leaseRsp.ID))

	return err
}
