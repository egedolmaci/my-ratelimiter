package ratelimiter

import (
	"testing"
	"time"
)

func TestIsRequestAllowed(t *testing.T) {
	t.Run("single request", func(t *testing.T) {
		rt := &Ratelimiter{limit: 1, storage: map[string]WindowData{}, windowSize: time.Minute}
		allowed, _ := rt.IsRequestAllowed("user123")

		if !allowed {
			t.Error("First request should be allowed")
		}

	})

	t.Run("2 requests at once", func(t *testing.T) {
		rt := &Ratelimiter{limit: 1, storage: map[string]WindowData{}, windowSize: time.Minute}
		rt.IsRequestAllowed("user123")
		allowed, _ := rt.IsRequestAllowed("user123")

		if allowed {
			t.Error("Second request should be disallowed when limit is 1")
		}
	})

	t.Run("1 request fills limit and limit opens up for another", func(t *testing.T) {
		rt := NewRateLimiter(1, time.Millisecond*100)
		defer rt.Stop()
		rt.IsRequestAllowed("user123")

		time.Sleep(150 * time.Millisecond)

		allowed, _ := rt.IsRequestAllowed("user123")

		if !allowed {
			t.Errorf("It must be allowed since windowSize amount of time has passed")
		}
	})
	t.Run("1 request fills the limit and since enough time has not passed limit is full", func(t *testing.T) {
		rt := NewRateLimiter(1, time.Millisecond*100)
		defer rt.Stop()
		rt.IsRequestAllowed("user123")

		time.Sleep(10 * time.Millisecond)

		allowed, _ := rt.IsRequestAllowed("user123")

		if allowed {
			t.Errorf("It must not be allowed since required time for the limit to reset has not passed")
		}
	})

}

func TestNewRatelimiter(t *testing.T) {
	rt := NewRateLimiter(10, time.Minute)
	defer rt.Stop()

	if rt.limit != 10 {
		t.Errorf("Expected limit 10, got %d", rt.limit)
	}

	if rt.storage == nil {
		t.Errorf("Storage should be initialized")
	}

	allowed, _ := rt.IsRequestAllowed("user123")

	if !allowed {
		t.Errorf("first request should be allowed")
	}
}
