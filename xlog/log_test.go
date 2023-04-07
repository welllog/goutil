package xlog

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func initTestLogger() {
	SetTimeFormat("")
	SetColor(false)
	SetCaller(false)
}

func TestPrintf(t *testing.T) {
	initTestLogger()

	tests := []struct {
		name string
		fn   func(string, ...any)
		lv   Level
	}{
		{
			name: "Errorf",
			fn:   Errorf,
			lv:   ERROR,
		},
		{
			name: "Warnf",
			fn:   Warnf,
			lv:   WARN,
		},
		{
			name: "Infof",
			fn:   Infof,
			lv:   INFO,
		},
		{
			name: "Debugf",
			fn:   Debugf,
			lv:   DEBUG,
		},
	}

	var buf bytes.Buffer
	SetWriter(NewWriter(&buf))

	for _, tt := range tests {
		SetEncode(PLAIN)
		tt.fn("test %s", "printf")
		want := fmt.Sprintf("\t%s\t%s\n", tt.lv.String(), "test printf")
		if buf.String() != want {
			t.Errorf("%s() = %s, want = %s", tt.name, buf.String(), want)
		}
		buf.Reset()

		SetEncode(JSON)
		tt.fn("test %s", "printf")
		want = fmt.Sprintf(`{"%s":"","%s":"%s","%s":"%s"}`, fieldTime, fieldLevel, tt.lv.String(), fieldContent, "test printf") + "\n"
		if buf.String() != want {
			t.Errorf("%s() = %s, want = %s", tt.name, buf.String(), want)
		}
		buf.Reset()
	}
}

func TestPrint(t *testing.T) {
	initTestLogger()

	tests := []struct {
		name string
		fn   func(...any)
		lv   Level
	}{
		{
			name: "Error",
			fn:   Error,
			lv:   ERROR,
		},
		{
			name: "Warn",
			fn:   Warn,
			lv:   WARN,
		},
		{
			name: "Info",
			fn:   Info,
			lv:   INFO,
		},
		{
			name: "Debug",
			fn:   Debug,
			lv:   DEBUG,
		},
	}

	var buf bytes.Buffer
	SetWriter(NewWriter(&buf))

	for _, tt := range tests {
		SetEncode(PLAIN)
		tt.fn("test ", "print")
		want := fmt.Sprintf("\t%s\t%s\n", tt.lv.String(), "test print")
		if buf.String() != want {
			t.Errorf("%s() = %s, want = %s", tt.name, buf.String(), want)
		}
		buf.Reset()

		SetEncode(JSON)
		tt.fn("test ", "print")
		want = fmt.Sprintf(`{"%s":"","%s":"%s","%s":"%s"}`, fieldTime, fieldLevel, tt.lv.String(), fieldContent, "test print") + "\n"
		if buf.String() != want {
			t.Errorf("%s() = %s, want = %s", tt.name, buf.String(), want)
		}
		buf.Reset()
	}
}

func TestPrintw(t *testing.T) {
	initTestLogger()

	tests := []struct {
		name string
		fn   func(string, ...Field)
		lv   Level
	}{
		{
			name: "Errorw",
			fn:   Errorw,
			lv:   ERROR,
		},
		{
			name: "Warnw",
			fn:   Warnw,
			lv:   WARN,
		},
		{
			name: "Infow",
			fn:   Infow,
			lv:   INFO,
		},
		{
			name: "Debugw",
			fn:   Debugw,
			lv:   DEBUG,
		},
	}

	var buf bytes.Buffer
	SetWriter(NewWriter(&buf))

	for _, tt := range tests {
		SetEncode(PLAIN)
		tt.fn("test", Field{Key: "age", Value: 18}, Field{Key: "addr", Value: "new york"})
		want := fmt.Sprintf("\t%s\t%s\t%s\t%s\t%s\t%s\n", tt.lv.String(), "test", "age", "18", "addr", "new york")
		if buf.String() != want {
			t.Errorf("%s() = %s, want = %s", tt.name, buf.String(), want)
		}
		buf.Reset()

		SetEncode(JSON)
		tt.fn("test", Field{Key: "age", Value: 18}, Field{Key: "addr", Value: "new york"})
		want = fmt.Sprintf(`{"%s":"","%s":"%s","%s":"%s","%s":%d,"%s":"%s"}`, fieldTime, fieldLevel, tt.lv.String(), fieldContent, "test", "age", 18, "addr", "new york") + "\n"
		if buf.String() != want {
			t.Errorf("%s() = %s, want = %s", tt.name, buf.String(), want)
		}
		buf.Reset()
	}
}

func TestSetLevel(t *testing.T) {
	var buf bytes.Buffer
	SetWriter(NewWriter(&buf))

	logging := func() {
		Error("test")
		Warn("test")
		Info("test")
		Debug("test")
	}

	getLines := func() int {
		var count int
		for {
			line, _ := buf.ReadBytes('\n')
			if len(line) == 0 {
				break
			}
			count++
		}
		return count
	}

	SetLevel(ERROR)
	logging()
	lines := getLines()
	want := 1
	if lines != want {
		t.Errorf("lines = %d, want = %d", lines, want)
	}

	buf.Reset()
	SetLevel(WARN)
	logging()
	lines = getLines()
	want = 2
	if lines != want {
		t.Errorf("lines = %d, want = %d", lines, want)
	}

	buf.Reset()
	SetLevel(INFO)
	logging()
	lines = getLines()
	want = 3
	if lines != want {
		t.Errorf("lines = %d, want = %d", lines, want)
	}

	buf.Reset()
	SetLevel(DEBUG)
	logging()
	lines = getLines()
	want = 4
	if lines != want {
		t.Errorf("lines = %d, want = %d", lines, want)
	}
}

func TestLogFacade(t *testing.T) {
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
	SetCaller(false)
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

	SetColor(false)
	for _, tt := range tests {
		logging(tt)
	}
}

func TestLog(t *testing.T) {
	SetEncode(PLAIN)
	Log("test log", WithLevel(INFO, "stat"), WithFields(Field{Key: "name", Value: "bob"}))
	Log("test log", WithLevel(WARN, "slow"), WithCaller(false))
	Log("test log", WithLevel(WARN, "slow"), WithCallerSkip(1), WithCallerSkip(-1))
}

type customLogger struct {
	Logger
}

func (l *customLogger) Slow(a ...any) {
	l.Log(fmt.Sprint(a...), WithLevel(WARN, "slow"), WithCallerSkipOne)
}

func (l *customLogger) Stat(a ...any) {
	Log(fmt.Sprint(a...), WithLevel(INFO, "stat"), WithCallerSkipOne)
}

func TestWrapLogger(t *testing.T) {
	SetEncode(PLAIN)
	l := customLogger{
		Logger: GetLogger(),
	}
	l.Slow("test slow")
	l.Stat("test stat")
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
