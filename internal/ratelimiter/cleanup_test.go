package ratelimiter

import (
	"fmt"
	"testing"
	"time"
    "github.com/egedolmaci/my-ratelimiter/internal/strategies"
)

func TestCleanup(t *testing.T) {
    strategy := strategies.NewFixedWindowStrategy(1, time.Microsecond * 100)
    rl := NewRateLimiterWithStrategy(strategy)
    defer rl.Stop()

    rl.strategy.IsRequestAllowed("user1")

    if strategy.GetStorageSize() != 1 {
        t.Error("Should have 1 entry")
    }

    time.Sleep(time.Millisecond * 250)

    if strategy.GetStorageSize() != 0 {
        t.Error("Should be cleaned up")
    }
}


func TestCleanupDoesNotAffectActiveRequests(t *testing.T) {
    strategy := strategies.NewFixedWindowStrategy(1, time.Microsecond * 100)
    rl := NewRateLimiterWithStrategy(strategy)
    defer rl.Stop()

	for i := 0; i < 10; i++ {
    	strategy.IsRequestAllowed(fmt.Sprintf("old-user-%d", i))
	}

    time.Sleep(time.Millisecond * 250)

	strategy.IsRequestAllowed("new-user-ege")

    if strategy.GetStorageSize() != 1 {
        t.Error("Should have 1 entry")
    }
}