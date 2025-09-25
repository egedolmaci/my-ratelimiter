package strategies

import (
	"testing"
	"time"
)

func TestIsRequestAllowedSlidingWindowLog(t *testing.T) {
	t.Run("should allow first request when limit is 1", func(t *testing.T) {
		rl := NewSlidingWindowLogStrategy(1, time.Minute)
		defer rl.Stop()

		allowed, _ := rl.IsRequestAllowed("ege")

		if !allowed {
			t.Errorf("First request must be allowed when limit is 1 got %t expected %t", allowed, true)
		}
	}) 

	t.Run("should not allow second request when limit is 1", func(t *testing.T) {
		rl := NewSlidingWindowLogStrategy(1, time.Minute)
		defer rl.Stop()

		rl.IsRequestAllowed("ege")
		allowed, _ := rl.IsRequestAllowed("ege")

		if allowed {
			t.Errorf("Second request must not be allowed when limit is 1")
		}
	})
}