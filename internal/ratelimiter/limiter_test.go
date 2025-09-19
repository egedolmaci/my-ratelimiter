package ratelimiter

import (
	"testing"
	"time"
)

func TestIsRequestAllowed(t *testing.T) {
	t.Run("single request", func(t *testing.T) {
		rt := &Ratelimiter{limit: 1, storage: map[string]WindowData{}}
		allowed := rt.IsRequestAllowed("user123")

		if !allowed {
			t.Error("First request should be allowed")
		}

	})

	t.Run("2 requests at once", func(t *testing.T) {
		rt := &Ratelimiter{limit: 1, storage: map[string]WindowData{}}
		rt.IsRequestAllowed("user123")
		allowed := rt.IsRequestAllowed("user123")

		if allowed {
			t.Error("Second request should be disallowed when limit is 1")
		}
	})

}

func TestNewRatelimiter(t *testing.T) {
	rt := NewRateLimiter(10, time.Minute)

	if rt.limit != 10 {
		t.Errorf("Expected limit 10, got %d", rt.limit)
	}

	if rt.storage == nil {
		t.Errorf("Storage should be initialized")
	}

	allowed := rt.IsRequestAllowed("user123")

	if !allowed {
		t.Errorf("first request should be allowed")
	}
}