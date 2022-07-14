package xgrpc

import (
	"context"
	"net/http"

	"google.golang.org/grpc/codes"
)

type Handler func(ctx context.Context, req *http.Request, writer ResponseWriter) error

type MiddlewareFunc func(ctx context.Context, req *http.Request, writer ResponseWriter, next Handler) error

type Middlewares struct {
	middlewares []MiddlewareFunc
}

func (m *Middlewares) Use(mwf ...MiddlewareFunc) {
	m.middlewares = append(m.middlewares, mwf...)
}

func (m *Middlewares) WrapHandler(handler http.Handler) http.Handler {
	max := len(m.middlewares)
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()

		var (
			chainHandler   Handler
			index, current int
		)

		cusWriter, ok := writer.(ResponseWriter)
		if !ok {
			cusWriter = NewResponseWriter(writer)
		}
		chainHandler = func(ctx context.Context, req *http.Request, writer ResponseWriter) error {
			current = index
			if current < max {
				index++
				return m.middlewares[current](ctx, req, writer, chainHandler)
			}

			req = req.WithContext(ctx)
			handler.ServeHTTP(writer, req)

			if code, msg := writer.GetBusinessErr(); code != 0 {
				return NewError(code, msg, writer.Status())
			} else if writer.Status() > 299 {
				return NewError(int(codes.Internal), http.StatusText(writer.Status()), writer.Status())
			}
			return nil
		}

		err := chainHandler(ctx, request, cusWriter)
		if !cusWriter.Written() && err != nil {
			// 处理中间件本身抛出的错误
			httpErr := Convert(err)
			_ = ErrResponse(httpErr, cusWriter)
		}
	})
}

func (m *Middlewares) Middleware() func(http.Handler) http.Handler {
	return m.WrapHandler
}
