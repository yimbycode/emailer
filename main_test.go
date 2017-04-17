package main

import (
	"net/http/httptest"
	"testing"

	google "github.com/kevinburke/google-oauth-handler"
)

func TestServerRedirects(t *testing.T) {
	mux := NewServeMux(google.NewAuthenticator(google.Config{
		SecretKey: NewRandomKey(),
	}), nil, "")
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != 302 {
		t.Errorf("GET /: got code %d, want 302", w.Code)
	}
}
