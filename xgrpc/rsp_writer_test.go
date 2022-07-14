package xgrpc

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewResponseWriter(t *testing.T) {
	recorder := httptest.NewRecorder()
	rsp := NewResponseWriter(recorder)
	rsp.WriteHeader(http.StatusBadGateway)

	if rsp.Status() != http.StatusBadGateway {
		t.Errorf("Status() = %d, want %d", rsp.Status(), http.StatusBadGateway)
	}
	if recorder.Code != http.StatusBadGateway {
		t.Errorf("recorder.Code = %d, want %d", recorder.Code, http.StatusBadGateway)
	}

	body := []byte("Hello World")
	rsp.Write(body)
	if recorder.Body.String() != string(body) {
		t.Errorf("recorder.Body = %s, want %s", recorder.Body.String(), string(body))
	}
	if rsp.Size() != len(body) {
		t.Errorf("Size() = %d, want %d", rsp.Size(), len(body))
	}
	if !rsp.Written() {
		t.Errorf("Written() = false, want true")
	}

	code, msg := 1, "test"
	rsp.setBusinessErr(code, msg)
	code2, msg2 := rsp.GetBusinessErr()
	if code2 != code || msg2 != msg {
		t.Errorf("GetBusinessErr() = (%d, %s), want (%d, %s)", code2, msg2, code, msg)
	}
}

func TestResponseWriter_Written(t *testing.T) {
	recorder := httptest.NewRecorder()
	rsp := NewResponseWriter(recorder)
	rsp.WriteHeader(http.StatusBadGateway)

	if !rsp.Written() {
		t.Errorf("Written() = false, want true")
	}
}
