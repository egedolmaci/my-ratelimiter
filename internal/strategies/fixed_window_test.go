package strategies

import (
	"fmt"
	"testing"
	"time"
)

func TestFixedWindowStrategy(t *testing.T) {
	t.Run("single request", func(t *testing.T) {
  		rt := NewFixedWindowStrategy(1, 100 * time.Millisecond)
		allowed, _ := rt.IsRequestAllowed("user123")
	
		if !allowed {
			t.Error("First request should be allowed")
		}
	
	})
	
	t.Run("2 requests at once", func(t *testing.T) {
  		rt := NewFixedWindowStrategy(1, 100 * time.Millisecond)
		rt.IsRequestAllowed("user123")
		allowed, _ := rt.IsRequestAllowed("user123")
	
		if allowed {
			t.Error("Second request should be disallowed when limit is 1")
		}
	})

}

func TestCleanup(t *testing.T) {
    strategy := NewFixedWindowStrategy(1, 100 * time.Millisecond)
    defer strategy.Stop()

    strategy.IsRequestAllowed("user1")

    if strategy.getStorageSize() != 1 {
        t.Error("Should have 1 entry")
    }

    time.Sleep(time.Millisecond * 250)

    if strategy.getStorageSize() != 0 {
        t.Error("Should be cleaned up")
    }
}

func TestCleanupDoesNotAffectActiveRequests(t *testing.T) {
    strategy := NewFixedWindowStrategy(1, 100 * time.Millisecond)
    defer strategy.Stop()

	for i := 0; i < 10; i++ {
    	strategy.IsRequestAllowed(fmt.Sprintf("old-user-%d", i))
	}

    time.Sleep(time.Millisecond * 250)

	strategy.IsRequestAllowed("new-user-ege")

    if strategy.getStorageSize() != 1 {
        t.Error("Should have 1 entry")
    }
}
