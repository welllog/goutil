package xgrpc

import (
	"bytes"
	"context"
	"net/http"

	"google.golang.org/grpc/status"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/grpc/codes"
)

var (
	_rspPrefix = []byte(`{"code":0,"msg":"success","data":`)
	_rspSuffix = []byte(`}`)
)

type GrpcStatus interface {
	GRPCStatus() *status.Status
}

type JsonBodyMarshaler struct {
	runtime.Marshaler
}

func (h *JsonBodyMarshaler) ContentType(v interface{}) string {
	if httpBody, ok := v.(*httpbody.HttpBody); ok {
		return httpBody.GetContentType()
	}
	return h.Marshaler.ContentType(v)
}

func (h *JsonBodyMarshaler) Marshal(v interface{}) ([]byte, error) {
	if httpBody, ok := v.(*httpbody.HttpBody); ok {
		return httpBody.Data, nil
	}
	buf := bytes.NewBuffer(_rspPrefix)
	err := h.NewEncoder(buf).Encode(v)
	if err != nil {
		return nil, err
	}
	buf.Write(_rspSuffix)
	return buf.Bytes(), nil
}

var (
	_fallback        = []byte(`{"code": 13, "msg": "failed to handle error", "data": {}}`)
	_routingFallback = []byte(`{"code": 13, "msg": "failed to handle routine error", "data": {}}`)
)

func HTTPErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler,
	w http.ResponseWriter, r *http.Request, err error) {

	httpErr := Convert(err)
	if err := ErrResponse(httpErr, w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(_fallback)
	}
}

func RoutingErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, httpStatus int) {
	sterr := NewError(int(codes.Internal), "Unexpected routing error", http.StatusInternalServerError)
	switch httpStatus {
	case http.StatusBadRequest:
		sterr = NewError(int(codes.InvalidArgument), http.StatusText(httpStatus), httpStatus)
	case http.StatusMethodNotAllowed:
		sterr = NewError(int(codes.Unimplemented), http.StatusText(httpStatus), httpStatus)
	case http.StatusNotFound:
		sterr = NewError(int(codes.NotFound), http.StatusText(httpStatus), httpStatus)
	}
	if err := ErrResponse(sterr, w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(_routingFallback)
	}
}
