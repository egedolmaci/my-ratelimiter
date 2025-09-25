package ratelimiter

import (
	"testing"
	"time"

	"github.com/egedolmaci/my-ratelimiter/internal/strategies"
)
func TestIsRequestAllowedFixedWindow(t *testing.T) {
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


