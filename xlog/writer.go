package xlog

import (
	"io"
	"os"
)

var csWriter = NewConsoleWriter()

type Writer interface {
	Write(level Level, p []byte) (n int, err error)
}

type consoleWriter struct {
	sw io.Writer
	ew io.Writer
}

func (c *consoleWriter) Write(level Level, p []byte) (n int, err error) {
	if level >= WARN {
		return c.ew.Write(p)
	}
	return c.sw.Write(p)
}

type customWriter struct {
	w io.Writer
}

func (c *customWriter) Write(level Level, p []byte) (n int, err error) {
	return c.w.Write(p)
}

func NewConsoleWriter() Writer {
	return &consoleWriter{
		sw: os.Stdout,
		ew: os.Stderr,
	}
}

func NewWriter(w io.Writer) Writer {
	return &customWriter{
		w: w,
	}
}
