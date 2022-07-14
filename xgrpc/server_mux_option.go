package xgrpc

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func ServerMuxOptions(ops ...runtime.ServeMuxOption) []runtime.ServeMuxOption {
	def := []runtime.ServeMuxOption{
		// default runtime.HTTPBodyMarshaler
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &JsonBodyMarshaler{
			Marshaler: _jsonMarshaler,
		}),
		// default runtime.DefaultHTTPErrorHandler
		runtime.WithErrorHandler(HTTPErrorHandler),
		runtime.WithStreamErrorHandler(runtime.DefaultStreamErrorHandler),
		// default runtime.DefaultRoutingErrorHandler
		runtime.WithRoutingErrorHandler(RoutingErrorHandler),
		runtime.WithIncomingHeaderMatcher(func(s string) (string, bool) {
			return "", false
		}),
	}
	return append(def, ops...)
}
