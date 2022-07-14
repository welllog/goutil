package xgrpc

import (
	"bytes"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoregistry"
)

var _jsonMarshaler runtime.Marshaler = &runtime.JSONPb{
	MarshalOptions: protojson.MarshalOptions{
		UseProtoNames:   true, // 默认false 使用camelCase
		EmitUnpopulated: true, // 默认true,为所有字段设置默认值
		Resolver:        protoregistry.GlobalTypes,
	},
	UnmarshalOptions: protojson.UnmarshalOptions{
		DiscardUnknown: true, // false 代表禁用未知字段
	},
}

var _errMarshalHandle = func(err *Error, m runtime.Marshaler) ([]byte, error) {
	buf := bytes.NewBuffer([]byte(`{"code":`))
	buf.WriteString(strconv.Itoa(err.errCode))
	buf.WriteString(`,"msg":"`)
	buf.WriteString(err.errMsg)
	buf.WriteString(`","data":`)
	if err.data != nil {
		if err := m.NewEncoder(buf).Encode(err.data); err != nil {
			return nil, err
		}
	} else {
		buf.WriteString("{}")
	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}

// should called be first
func WithDefaultMarshaler(m runtime.Marshaler) {
	_jsonMarshaler = m
}

func WithErrMarshalHandle(fn func(err *Error, m runtime.Marshaler) ([]byte, error)) {
	_errMarshalHandle = fn
}
