package xlog

import (
	"context"
	"testing"
)

func TestWithContext(t *testing.T) {
	setDefLogger(newLogger())

	ctx := context.WithValue(context.Background(), "uid", 3)
	SetDefCtxHandle(func(ctx context.Context) []Field {
		var fs []Field
		uid, ok := ctx.Value("uid").(int)
		if ok {
			fs = append(fs, Field{Key: "uid", Value: uid})
		}

		name, ok := ctx.Value("name").(string)
		if ok {
			fs = append(fs, Field{Key: "name", Value: name})
		}
		return fs
	})

	SetEncode(PLAIN)
	l := WithContext(GetLogger(), ctx)
	l.Debug("test uid")

	ctx = context.WithValue(context.Background(), "name", "bob")
	l = WithContext(l, ctx)
	l.Debug("test name")

	l = WithContext(l, context.Background())
	l.Debug("test name")

	l = WithEntries(l, map[string]any{
		"ip":      "127.0.0.1",
		"score":   99.9,
		"success": true,
	})

	l.Debug("test entries")

	l = WithContext(l, context.WithValue(context.Background(), "name", "linda"))
	l.Debug("test override name")
	l.Log("test log", WithLevel(DEBUG, "trace"))
}
