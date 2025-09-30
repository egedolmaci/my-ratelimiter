package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/egedolmaci/my-ratelimiter/internal/ratelimiter"
	"github.com/egedolmaci/my-ratelimiter/internal/strategies"
)

func TestRateLimitMiddleware(t *testing.T) {
	t.Run("single request", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		}

		rl := ratelimiter.NewRateLimiter(1, time.Minute, &strategies.RealTimeProvider{})
		defer rl.Stop()
		middleware := Middleware{Ratelimiter: rl}

		next := middleware.RateLimitMiddleware(handler)

		req := httptest.NewRequest("GET", "/test", nil)
		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", rec.Code)
		}

	})
	
	t.Run("2 requests when limit is one", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		}

		rl := ratelimiter.NewRateLimiter(1, time.Minute, &strategies.RealTimeProvider{})
		defer rl.Stop()
		middleware := Middleware{Ratelimiter: rl}

		next := middleware.RateLimitMiddleware(handler)

		req := httptest.NewRequest("GET", "/test", nil)
		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, req)

		req2 := httptest.NewRequest("GET", "/test", nil)
		rec2 := httptest.NewRecorder()
		next.ServeHTTP(rec2, req2)

		if rec2.Code != http.StatusTooManyRequests {
			t.Errorf("Expected 429, got %d", rec.Code)
		}
	})
}


func TestMiddlewareLimit(t *testing.T) {
	t.Run("single request body", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}

		rl := ratelimiter.NewRateLimiter(10, time.Minute, &strategies.RealTimeProvider{})
		defer rl.Stop()
		middleware := Middleware{Ratelimiter: rl}

		next := middleware.RateLimitMiddleware(handler)

		req := httptest.NewRequest("GET", "/test", nil)
		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, req)

		expected := "Remaining limit = 9\n"
		got := rec.Body.String()

		if got != expected {
			t.Errorf("In the recorder body remaining limit -> want %s got %s", expected, got)
		}
	})
	t.Run("2 requests body", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}

		rl := ratelimiter.NewRateLimiter(10, time.Minute, &strategies.RealTimeProvider{})
		defer rl.Stop()
		middleware := Middleware{Ratelimiter: rl}

		next := middleware.RateLimitMiddleware(handler)

		req := httptest.NewRequest("GET", "/test", nil)
		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, req)

		req2 := httptest.NewRequest("GET", "/test", nil)
		rec2 := httptest.NewRecorder()
		next.ServeHTTP(rec2, req2)

		expected := "Remaining limit = 8\n"
		got := rec2.Body.String()

		if got != expected {
			t.Errorf("In the recorder body remaining limit -> want %s got %s", expected, got)
		}
	})
}