package xgrpc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"runtime"
	"strings"

	"google.golang.org/grpc/codes"
)

const (
	_MAX_SIZE = 1 << 20
)

type LevelLogger interface {
	Debugw(msg string, fields ...interface{})
	Infow(msg string, fields ...interface{})
	Warnw(msg string, fields ...interface{})
	Errw(msg string, fields ...interface{})
}

type codeErr interface {
	Code() codes.Code
}

func LoggingAndRecover(logger LevelLogger, debug bool) MiddlewareFunc {
	return func(ctx context.Context, req *http.Request, writer ResponseWriter, next Handler) (err error) {
		defer func() {
			if p := recover(); p != nil {
				pc, _, no, _ := runtime.Caller(4) // 0 此处   1 runtime.Panic
				var buf [2048]byte
				n := runtime.Stack(buf[:], false)

				panicInfo := fmt.Sprintf(
					"panic:%v,callee:%s,on line %d;stack:%s",
					p, runtime.FuncForPC(pc).Name(), no, string(buf[:n]),
				)

				logPanic(ctx, logger, req, panicInfo)
				err = NewError(http.StatusInternalServerError, "server error", http.StatusInternalServerError)
				return
			}
		}()

		var body string
		if debug { // debug时，记录请求body
			body, err = getBody(req)
			if err != nil {
				return
			}
		}

		err = next(ctx, req, writer)

		if debug {
			logDebug(ctx, logger, req, body, err)
			return
		}

		if err != nil {
			logErr(ctx, logger, req, err)
		}

		return
	}
}

func logPanic(ctx context.Context, logger LevelLogger, req *http.Request, panicInfo string) {
	fields := make([]interface{}, 0, 14)
	fields = fillFields(ctx, req, fields)
	fields = append(fields, "panic", panicInfo)
	logger.Errw("http.logging", fields...)
}

func logDebug(ctx context.Context, logger LevelLogger, req *http.Request, body string, err error) {
	fields := make([]interface{}, 0, 16)
	fields = fillFields(ctx, req, fields)
	fields = append(fields, "body", body)
	if err != nil {
		handleErr(logger, fields, err)
		return
	}
	logger.Debugw("http.logging", fields...)
}

func logErr(ctx context.Context, logger LevelLogger, req *http.Request, err error) {
	fields := make([]interface{}, 0, 14)
	fields = fillFields(ctx, req, fields)
	handleErr(logger, fields, err)
}

func fillFields(ctx context.Context, req *http.Request, fields []interface{}) []interface{} {
	fields = append(fields, "uri", req.URL.Path)
	fields = append(fields, "query", req.URL.RawQuery)
	return fields
}

func handleErr(logger LevelLogger, fields []interface{}, err error) {
	fields = append(fields, "err", err.Error())
	if errors.Is(err, context.Canceled) { // 该错误标记为info
		logger.Infow("http.logging", fields...)
		return
	}

	if cerr, ok := err.(codeErr); ok {
		ecode := cerr.Code()
		if ecode == codes.Unknown || ecode == codes.Internal { // 内部错误，记录警告日志
			logger.Warnw("http.logging", fields...)
			return
		}
		// 普通业务错误
		logger.Infow("http.logging", fields...)
		return
	}

	// 内部错误
	logger.Warnw("http.logging", fields...)
}

type bytesReader interface {
	Bytes() []byte
}

func getBody(req *http.Request) (body string, err error) {
	switch req.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
	default:
		return
	}
	ct := req.Header.Get("Content-Type")
	if ct == "" {
		return
	}
	ct, _, _ = mime.ParseMediaType(ct)
	if ct == "application/json" || ct == "application/x-www-form-urlencoded" {
		if sr, ok := req.Body.(bytesReader); ok {
			body = strings.ReplaceAll(string(sr.Bytes()), "\"", "'") // 替换双引号，防止日志抓取错误
			return
		}
		var b []byte
		b, err = io.ReadAll(io.LimitReader(req.Body, _MAX_SIZE))
		if err != nil {
			return
		}
		_ = req.Body.Close()
		body = strings.ReplaceAll(string(b), "\"", "'") // 替换双引号，防止日志抓取错误
		req.Body = io.NopCloser(bytes.NewBuffer(b))
	}
	return
}
