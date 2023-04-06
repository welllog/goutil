package xlog

import (
	"fmt"
	"os"
	"testing"
)

func TestGetLogger(t *testing.T) {
	tests := []struct {
		arr    []any
		f      string
		a      []any
		msg    string
		fields []Field
	}{
		{
			arr: []any{"t1", "t2"},
			f:   "t%d%s",
			a:   []any{3, "t4"},
			msg: "t5",
			fields: []Field{
				{Key: "name", Value: "bob"},
				{Key: "age", Value: 18},
			},
		},
	}

	for _, tt := range tests {
		logging(tt)
	}

	SetEncode(PLAIN)
	for _, tt := range tests {
		logging(tt)
	}

	SetLevel(WARN)
	for _, tt := range tests {
		logging(tt)
	}

	SetEnableCaller(false)
	for _, tt := range tests {
		logging(tt)
	}

	SetTimeFormat("2006-01-02 15:04:05")
	for _, tt := range tests {
		logging(tt)
	}

	f, err := os.Create("test.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	SetWriter(NewWriter(f))
	SetTimeFormat("")
	for _, tt := range tests {
		logging(tt)
	}

	SetEncode(JSON)
	for _, tt := range tests {
		logging(tt)
	}
}

func TestFatal(t *testing.T) {
	Fatal("t1", "t2")
	Fatal("t1", "t2", "t3")
}

func TestLogger_Log(t *testing.T) {
	SetEncode(PLAIN)
	GetLogger().Log("t1", WithLevel(INFO, "stat"))
	GetLogger().Log("t1", WithLevel(WARN, "slow"))
}

type customLogger struct {
	Logger
}

func (l *customLogger) Slow(a ...any) {
	l.Log(fmt.Sprint(a...), WithLevel(WARN, "slow"), WithCallerSkipOne)
}

func TestWrapLogger(t *testing.T) {
	SetEncode(PLAIN)
	l := customLogger{
		Logger: GetLogger(),
	}
	l.Slow("t1", "t2")
	l.Slow("test slow")
}

func logging(tt struct {
	arr    []any
	f      string
	a      []any
	msg    string
	fields []Field
}) {
	Debug(tt.arr...)
	Debugf(tt.f, tt.a...)
	Debugw(tt.msg, tt.fields...)

	Info(tt.arr...)
	Infof(tt.f, tt.a...)
	Infow(tt.msg, tt.fields...)

	Warn(tt.arr...)
	Warnf(tt.f, tt.a...)
	Warnw(tt.msg, tt.fields...)

	Error(tt.arr...)
	Errorf(tt.f, tt.a...)
	Errorw(tt.msg, tt.fields...)
}
