package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJsonMiddleware(t *testing.T) {
	t.Run("sets content-type header", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("test"))
		})

		middleware := JsonMiddleware(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		middleware.ServeHTTP(w, req)

		contentType := w.Header().Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("expected Content-Type 'application/json', got '%s'", contentType)
		}

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		body := w.Body.String()
		if body != "test" {
			t.Errorf("expected body 'test', got '%s'", body)
		}
	})

	t.Run("calls next handler", func(t *testing.T) {
		called := false
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
		})

		middleware := JsonMiddleware(handler)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		middleware.ServeHTTP(w, req)

		if !called {
			t.Error("expected next handler to be called")
		}
	})

	t.Run("preserves other operations", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Custom-Header", "custom-value")
			w.WriteHeader(http.StatusCreated)
		})

		middleware := JsonMiddleware(handler)

		req := httptest.NewRequest(http.MethodPost, "/", nil)
		w := httptest.NewRecorder()

		middleware.ServeHTTP(w, req)

		if w.Header().Get("Content-Type") != "application/json" {
			t.Error("Content-Type not set")
		}

		if w.Header().Get("X-Custom-Header") != "custom-value" {
			t.Error("Custom header not preserved")
		}

		if w.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", w.Code)
		}
	})
}
