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
}