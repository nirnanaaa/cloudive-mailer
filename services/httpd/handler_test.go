package httpd_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func GetHttpHandler(rr *httptest.ResponseRecorder, req *http.Request) error {
	httpx := CreateService(false)
	handler := http.HandlerFunc(httpx.Handler.ServeHTTP)
	handler.ServeHTTP(rr, req)
	return nil
}

// Ensure the handler handles status requests correctly.
func TestHandler_Status(t *testing.T) {
	w := httptest.NewRecorder()
	req := MustNewRequest("GET", "/healthz", nil)
	GetHttpHandler(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: %d", w.Code)
	}
	GetHttpHandler(w, MustNewRequest("HEAD", "/healthz", nil))
	if w.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: %d", w.Code)
	}
}

// MustNewRequest returns a new HTTP request. Panic on error.
func MustNewRequest(method, urlStr string, body io.Reader) *http.Request {
	r, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		panic(err.Error())
	}
	return r
}
