package ratelimiter

import (
	"testing"
	"time"

	"github.com/egedolmaci/my-ratelimiter/internal/strategies"
)
func TestIsRequestAllowedFixedWindow(t *testing.T) {
	t.Run("single request", func(t *testing.T) {
		rt := NewRateLimiterWithStrategy(strategies.NewFixedWindowStrategy(1, time.Millisecond*100))
		allowed, _ := rt.IsRequestAllowed("user123")

		if !allowed {
			t.Error("First request should be allowed")
		}

	})

	t.Run("2 requests at once", func(t *testing.T) {
		rt := NewRateLimiterWithStrategy(strategies.NewFixedWindowStrategy(1, time.Millisecond*100))
		rt.IsRequestAllowed("user123")
		allowed, _ := rt.IsRequestAllowed("user123")

		if allowed {
			t.Error("Second request should be disallowed when limit is 1")
		}
	})

	t.Run("1 request fills limit and limit opens up for another", func(t *testing.T) {
		rt := NewRateLimiterWithStrategy(strategies.NewFixedWindowStrategy(1, time.Millisecond*100))
		defer rt.Stop()
		rt.IsRequestAllowed("user123")

		time.Sleep(150 * time.Millisecond)

		allowed, _ := rt.IsRequestAllowed("user123")

		if !allowed {
			t.Errorf("It must be allowed since windowSize amount of time has passed")
		}
	})
	t.Run("1 request fills the limit and since enough time has not passed limit is full", func(t *testing.T) {
		rt := NewRateLimiterWithStrategy(strategies.NewFixedWindowStrategy(1, time.Millisecond*100))
		defer rt.Stop()
		rt.IsRequestAllowed("user123")

		time.Sleep(10 * time.Millisecond)

		allowed, _ := rt.IsRequestAllowed("user123")

		if allowed {
			t.Errorf("It must not be allowed since required time for the limit to reset has not passed")
		}
	})
}

func TestIsRequestAllowedSlidingWindowLog(t *testing.T) {
	t.Run("should allow first request when limit is 1", func(t *testing.T) {
		rl := NewRateLimiterWithStrategy(strategies.NewSlidingWindowLogStrategy(1, time.Minute))
		defer rl.Stop()

		allowed, _ := rl.IsRequestAllowed("ege")

		if !allowed {
			t.Errorf("First request must be allowed when limit is 1 got %t expected %t", allowed, true)
		}
	}) 
}
