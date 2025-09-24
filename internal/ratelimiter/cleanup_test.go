package ratelimiter

import (
	"testing"
	"time"

	"github.com/egedolmaci/my-ratelimiter/internal/strategies"
)

func TestStop(t *testing.T) {
    rl := NewRateLimiterWithStrategy(strategies.NewFixedWindowStrategy(1, time.Microsecond * 100))
    done := make(chan bool, 1)

    go func() {
        rl.Stop()
        done <- true
    }()

    select {
    case <-done:
    case <-time.After(1 * time.Second):
        t.Fatal("Stop() hung or took too long")
    }
}


func TestStopDoesNotAffectActiveRequests(t *testing.T) {
    rl := NewRateLimiterWithStrategy(strategies.NewFixedWindowStrategy(1, time.Microsecond * 100))

	allowed, _ := rl.IsRequestAllowed("new-user-ege")
    rl.Stop()

    if !allowed {
        t.Errorf("Stop() should not affect active requests request should be allowed, got %t expected %t", allowed, true)
    }
}