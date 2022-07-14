package xgrpc

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const contentType = "application/json; charset=utf-8"

var _ error = (*Error)(nil)

type Error struct {
	err      error
	httpCode int
	errCode  int
	errMsg   string
	data     proto.Message
}

func NewError(errCode int, errMsg string, httpCode ...int) *Error {
	return WrapError(nil, errCode, errMsg, httpCode...)
}

func WrapError(err error, errCode int, errMsg string, httpCode ...int) *Error {
	code := http.StatusOK
	if len(httpCode) > 0 {
		code = httpCode[0]
	}
	return &Error{
		err:      err,
		httpCode: code,
		errCode:  errCode,
		errMsg:   errMsg,
	}
}

func WrapGrpcStatus(st *status.Status, errCode int, errMsg string, httpCode ...int) *Error {
	e := WrapError(st.Err(), errCode, errMsg, httpCode...)
	spb := st.Proto()
	if len(spb.Details) > 0 {
		return e.WithData(spb.Details[0])
	}
	return e
}

func (h *Error) WithData(data proto.Message) *Error {
	h.data = data
	return h
}

func (h *Error) Is(err error) bool {
	if h.err == nil {
		return false
	}
	return errors.Is(h.err, err)
}

func (h *Error) Error() string {
	msg := "[" + strconv.Itoa(h.errCode) + "]" + h.errMsg
	if h.err != nil {
		return msg + "; wrapped: " + h.err.Error()
	}
	return msg
}

func (h *Error) Unwrap() error {
	return h.err
}

func (h *Error) StatusCode() int {
	return h.httpCode
}

func (h *Error) Code() codes.Code {
	return codes.Code(h.errCode)
}

func (h *Error) GRPCStatus() *status.Status {
	st := status.New(codes.Code(h.errCode), h.errMsg)
	if h.data != nil {
		nst, err := st.WithDetails(h.data)
		if err == nil {
			return nst
		}
	}
	return st
}

func Convert(err error) *Error {
	httpErr, ok := err.(*Error)
	if ok {
		return httpErr
	}

	if grpcStatus, ok := err.(GrpcStatus); ok {
		st := grpcStatus.GRPCStatus()
		code := st.Code()
		if code == codes.Unknown || code == codes.Internal {
			return WrapGrpcStatus(st, int(codes.Internal), "internal server error")
		}
		return WrapGrpcStatus(st, int(st.Code()), st.Message())
	}
	return WrapError(err, int(codes.Internal), "internal server error")
}

func ErrResponse(e *Error, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(e.httpCode)
	if nw, ok := w.(ResponseWriter); ok {
		errMsg := e.errMsg
		if e.err != nil {
			errMsg = e.err.Error()
		}
		nw.setBusinessErr(e.errCode, errMsg)
	}
	b, err := _errMarshalHandle(e, _jsonMarshaler)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}
