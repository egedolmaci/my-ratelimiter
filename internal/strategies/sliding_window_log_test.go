package strategies

import (
	"testing"
	"time"
)

func TestSlidingWindowLogStrategy(t *testing.T) {
	t.Run("single request", func(t *testing.T) {
		strategy := SlidingWindowLogStrategy{2, time.Minute}
	
		allowed, _ := strategy.IsRequestAllowed("ege")
	
		if !allowed {
			t.Errorf("Request should be allowed got %t, expected %t", allowed, true)
		}
	})

	t.Run("# of requests above limit", func(t *testing.T) {
		strategy := SlidingWindowLogStrategy{2, time.Minute}

		strategy.IsRequestAllowed("ege")
		strategy.IsRequestAllowed("ege")
		allowed, _ := strategy.IsRequestAllowed("ege")

		if allowed {
			t.Errorf("Limit is 2, this is 3rd request should not be allowed, got %t, expected %t", allowed, false)
		}
	})
	
}