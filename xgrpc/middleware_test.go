package xgrpc

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddlewares_Middleware(t *testing.T) {
	mids := &Middlewares{}
	max := 10
	for i := 0; i < max; i++ {
		mids.Use(func(ctx context.Context, req *http.Request, writer ResponseWriter, next Handler) error {
			key := "lv"
			value, ok := ctx.Value(key).(int)
			if !ok {
				value = 0
			}
			value++
			ctx = context.WithValue(ctx, key, value)
			return next(ctx, req, writer)
		})
	}

	handler := mids.WrapHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lv := r.Context().Value("lv").(int)
		if lv != max {
			t.Errorf("expected lv %d, got %d", max, lv)
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rsp := httptest.NewRecorder()
	handler.ServeHTTP(rsp, req)
}

func TestMiddlewares_Middleware2(t *testing.T) {
	mids := &Middlewares{}
	mids.Use(func(ctx context.Context, req *http.Request, writer ResponseWriter, next Handler) error {
		err := next(ctx, req, writer)
		if err == nil {
			t.Errorf("expected error")
		}
		t.Log(err.Error())
		return err
	})

	handler1 := mids.WrapHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))

	handler2 := mids.WrapHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wr, ok := w.(ResponseWriter)
		if ok {
			wr.setBusinessErr(1, "test")
		}
		w.WriteHeader(http.StatusBadGateway)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rsp := httptest.NewRecorder()
	handler1.ServeHTTP(rsp, req)

	handler2.ServeHTTP(rsp, req)
}

func TestMiddlewares_Middleware3(t *testing.T) {
	mids := &Middlewares{}
	mids.Use(func(ctx context.Context, req *http.Request, writer ResponseWriter, next Handler) error {
		err := next(ctx, req, writer)
		if err == nil {
			t.Errorf("expected error")
		}
		t.Log(err.Error())
		return err
	})
	mids.Use(func(ctx context.Context, req *http.Request, writer ResponseWriter, next Handler) error {
		return NewError(1, "test", http.StatusBadRequest)
	})

	handler := mids.WrapHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("expected abort")
	}))
	req := httptest.NewRequest("GET", "/", nil)
	rsp := httptest.NewRecorder()
	handler.ServeHTTP(rsp, req)

	if rsp.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, rsp.Code)
	}
}
