package ratelimiter

import (
	"fmt"
	"testing"
	"time"
)

func TestCleanup(t *testing.T) {
    strategy := NewFixedWindowStrategy(1, time.Millisecond * 100)
    defer strategy.Stop()

    strategy.IsRequestAllowed("user1")

    if strategy.GetStorageSize() != 1 {
        t.Error("Should have 1 entry")
    }

    time.Sleep(time.Millisecond * 250)

    if strategy.GetStorageSize() != 0 {
        t.Error("Should be cleaned up")
    }
}


func TestCleanupDoesNotAffectActiveRequests(t *testing.T) {
    strategy := NewFixedWindowStrategy(10, time.Millisecond * 100)
    defer strategy.Stop()

	for i := 0; i < 10; i++ {
    	strategy.IsRequestAllowed(fmt.Sprintf("old-user-%d", i))
	}

	

    time.Sleep(time.Millisecond * 250)


	strategy.IsRequestAllowed("new-user-ege")

    if strategy.GetStorageSize() != 1 {
        t.Error("Should have 1 entry")
    }
}