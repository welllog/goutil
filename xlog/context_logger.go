package xlog

import (
	"context"
	"fmt"
	"sync/atomic"
	"unsafe"
)

var defCtxHandler unsafe.Pointer

func init() {
	SetDefCtxHandle(emptyCtxHandle)
}

type CtxHandler struct {
	handle func(ctx context.Context) []Field
}

func SetDefCtxHandle(handle func(ctx context.Context) []Field) {
	atomic.StorePointer(&defCtxHandler, unsafe.Pointer(&CtxHandler{handle: handle}))
}

func getDefCtxHandle() func(ctx context.Context) []Field {
	return (*CtxHandler)(atomic.LoadPointer(&defCtxHandler)).handle
}

func emptyCtxHandle(ctx context.Context) []Field {
	return nil
}

type contextLogger struct {
	Logger
	ctx       context.Context
	entries   map[string]any
	ctxHandle func(ctx context.Context) []Field
}

func WithContext(logger Logger, ctx context.Context, handles ...func(context.Context) []Field) Logger {
	var handle func(context.Context) []Field
	if len(handles) > 0 {
		handle = handles[0]
	} else {
		handle = getDefCtxHandle()
	}
	return contextLogger{
		Logger:    logger,
		ctx:       ctx,
		ctxHandle: handle,
	}
}

func WithEntries(logger Logger, entries map[string]any) Logger {
	return contextLogger{
		Logger:    logger,
		ctx:       context.Background(),
		entries:   entries,
		ctxHandle: emptyCtxHandle,
	}
}

func (c contextLogger) Log(content any, opts ...LogOption) {
	opts = append(opts, WithCallerSkipOne, WithFields(c.buildFields()...))
	c.Logger.Log(content, opts...)
}

func (c contextLogger) Fatal(a ...any) {
	if c.isOut(FATAL) {
		c.Logger.Log(fmt.Sprint(a...), withFatal, WithCallerSkipOne, WithFields(c.buildFields()...))
	}
}

func (c contextLogger) Fatalf(format string, a ...any) {
	if c.isOut(FATAL) {
		c.Logger.Log(fmt.Sprintf(format, a...), withFatal, WithCallerSkipOne, WithFields(c.buildFields()...))
	}
}

func (c contextLogger) Fatalw(msg string, fields ...Field) {
	if c.isOut(FATAL) {
		c.Logger.Log(msg, withFatal, WithCallerSkipOne, WithFields(c.buildFields(fields...)...))
	}
}

func (c contextLogger) Error(a ...any) {
	if c.isOut(ERROR) {
		c.Logger.Log(fmt.Sprint(a...), withError, WithCallerSkipOne, WithFields(c.buildFields()...))
	}
}

func (c contextLogger) Errorf(format string, a ...any) {
	if c.isOut(ERROR) {
		c.Logger.Log(fmt.Sprintf(format, a...), withError, WithCallerSkipOne, WithFields(c.buildFields()...))
	}
}

func (c contextLogger) Errorw(msg string, fields ...Field) {
	if c.isOut(ERROR) {
		c.Logger.Log(msg, withError, WithCallerSkipOne, WithFields(c.buildFields(fields...)...))
	}
}

func (c contextLogger) Warn(a ...any) {
	if c.isOut(WARN) {
		c.Logger.Log(fmt.Sprint(a...), withWarn, WithCallerSkipOne, WithFields(c.buildFields()...))
	}
}

func (c contextLogger) Warnf(format string, a ...any) {
	if c.isOut(WARN) {
		c.Logger.Log(fmt.Sprintf(format, a...), withWarn, WithCallerSkipOne, WithFields(c.buildFields()...))
	}
}

func (c contextLogger) Warnw(msg string, fields ...Field) {
	if c.isOut(WARN) {
		c.Logger.Log(msg, withWarn, WithCallerSkipOne, WithFields(c.buildFields(fields...)...))
	}
}

func (c contextLogger) Info(a ...any) {
	if c.isOut(INFO) {
		c.Logger.Log(fmt.Sprint(a...), withInfo, WithCallerSkipOne, WithFields(c.buildFields()...))
	}
}

func (c contextLogger) Infof(format string, a ...any) {
	if c.isOut(INFO) {
		c.Logger.Log(fmt.Sprintf(format, a...), withInfo, WithCallerSkipOne, WithFields(c.buildFields()...))
	}
}

func (c contextLogger) Infow(msg string, fields ...Field) {
	if c.isOut(INFO) {
		c.Logger.Log(msg, withInfo, WithCallerSkipOne, WithFields(c.buildFields(fields...)...))
	}
}

func (c contextLogger) Debug(a ...any) {
	if c.isOut(DEBUG) {
		c.Logger.Log(fmt.Sprint(a...), withDebug, WithCallerSkipOne, WithFields(c.buildFields()...))
	}
}

func (c contextLogger) Debugf(format string, a ...any) {
	if c.isOut(DEBUG) {
		c.Logger.Log(fmt.Sprintf(format, a...), withDebug, WithCallerSkipOne, WithFields(c.buildFields()...))
	}
}

func (c contextLogger) Debugw(msg string, fields ...Field) {
	if c.isOut(DEBUG) {
		c.Logger.Log(msg, withDebug, WithCallerSkipOne, WithFields(c.buildFields(fields...)...))
	}
}

func (c contextLogger) buildFields(fields ...Field) []Field {
	fields = append(fields, c.ctxHandle(c.ctx)...)
	if len(c.entries) > 0 {
		for key, entry := range c.entries {
			fields = append(fields, Field{Key: key, Value: entry})
		}
	}
	return fields
}

func withFatal(o *logOption) {
	o.level = FATAL
}

func withError(o *logOption) {
	o.level = ERROR
}

func withWarn(o *logOption) {
	o.level = WARN
}

func withInfo(o *logOption) {
	o.level = INFO
}

func withDebug(o *logOption) {
	o.level = DEBUG
}
