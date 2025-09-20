package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/egedolmaci/my-ratelimiter/internal/ratelimiter"
	"github.com/egedolmaci/my-ratelimiter/pkg/middleware"
)

func TestMain(t *testing.T) {
	t.Run("TestUnlimitedEndpoint", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/unlimited", nil)
		w := httptest.NewRecorder()

		UnlimitedHandler(w, req)

		if w.Code != 200 {
			t.Errorf("got %d, want %d", w.Code, 200)
		}
	})

	t.Run("TestLimitedEndpoint", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/limited", nil)
		w := httptest.NewRecorder()

		LimitedHandler(w, req)

		if w.Code != 429 {
			t.Errorf("got %d, want %d", w.Code, 429)
		}
	})

	t.Run("TestHandlersThroughServer", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/limited", nil)
		w := httptest.NewRecorder()

		req2 := httptest.NewRequest("GET", "/unlimited", nil)
		w2 := httptest.NewRecorder()

		server := NewServer()
		server.ServeHTTP(w, req)

		server.ServeHTTP(w2, req2)

		if w.Code != 429 {
			t.Errorf("for limited endpoint got %d, want %d", w.Code, 429)
		}

		if w2.Code != 200 {
			t.Errorf("for unlimited endpoint got %d, want %d", w2.Code, 200)
		}
	})

}

func TestIntegrationWithMiddleware(t *testing.T) {
	server := NewServer()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("rate limited"))

	})
	middleware := middleware.Middleware{Ratelimiter: ratelimiter.NewRateLimiter(2, time.Minute)}
	ratelimitedHandler := middleware.RateLimitMiddleware(handler)
	server.mux.HandleFunc("/test-ratelimited", ratelimitedHandler)

	t.Run("single request to server", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test-ratelimited", nil)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		if res.Code != http.StatusOK {
			t.Errorf("Status must be ok got %q", res.Code)
		}
	})

	t.Run("double request to server", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test-ratelimited", nil)
		res := httptest.NewRecorder()

		req2 := httptest.NewRequest("GET", "/test-ratelimited", nil)
		res2 := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		server.ServeHTTP(res2, req2)

		if res2.Code != http.StatusTooManyRequests{
			t.Errorf("Status must be 429 got %d", res.Code)
		}
	})

}

func TestIntegrationWithMiddlewareAdvanced(t *testing.T) {
		server := NewServer()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("rate limited"))

		})
		middleware := middleware.Middleware{Ratelimiter: ratelimiter.NewRateLimiter(3, time.Minute)}
		ratelimitedHandler := middleware.RateLimitMiddleware(handler)
		server.mux.HandleFunc("/test-ratelimited", ratelimitedHandler)

		t.Run("3 requests to server with limit 3 4th should fail", func(t *testing.T) {
			i := 0
			for {
				req := httptest.NewRequest("GET", "/test-ratelimited", nil)
				res := httptest.NewRecorder()

				server.ServeHTTP(res, req)
				i++

				if i == 3 {
					break
				}

			}
			req := httptest.NewRequest("GET", "/test-ratelimited", nil)
			res := httptest.NewRecorder()

			server.ServeHTTP(res, req)

			if res.Code != http.StatusTooManyRequests {
				t.Errorf("Status must be 429 got %q", res.Code)
			}
		})
}