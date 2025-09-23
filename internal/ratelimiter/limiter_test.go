package ratelimiter

import (
	"testing"
	"time"
)

func TestIsRequestAllowed(t *testing.T) {
	t.Run("single request", func(t *testing.T) {
		rt := &FixedWindowStrategy{limit: 1, storage: map[string]WindowData{}, windowSize: time.Minute}
		allowed, _ := rt.IsRequestAllowed("user123")

		if !allowed {
			t.Error("First request should be allowed")
		}

	})

	t.Run("2 requests at once", func(t *testing.T) {
		rt := &FixedWindowStrategy{limit: 1, storage: map[string]WindowData{}, windowSize: time.Minute}
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

func TestRateLimiterWithFixedWindowStrategy(t *testing.T) {
	t.Run("can use strategy-based rate limiter with same behaviour", func(t *testing.T) {
		strategy := NewFixedWindowStrategy(2, 100 * time.Millisecond)

		rl := NewRateLimiterWithStrategy(strategy)
		defer rl.Stop()

		allowed, remaining := rl.IsRequestAllowed("user1")
		if !allowed || remaining != 1 {
			t.Errorf("Expected allowed=true, remaining=1, got allowed=%v, remaining=%d", allowed, remaining) 
		}

		
		allowed, remaining = rl.IsRequestAllowed("user1")
		if !allowed || remaining != 0 {
			t.Errorf("Expected allowed=true, remaining=0, got allowed=%v, remaining=%d", allowed, remaining) 
		}

		allowed, remaining = rl.IsRequestAllowed("user1")
		if allowed || remaining != 0 {
			t.Errorf("Expected allowed=false, remaining=0, got allowed=%v, remaining=%d", allowed, remaining) 
		}
	})
}
