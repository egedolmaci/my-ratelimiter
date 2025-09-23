package ratelimiter

import (
	"fmt"
	"testing"
	"time"
)

func TestCleanup(t *testing.T) {
	rt := NewRateLimiter(1, time.Millisecond * 100)
	defer rt.Stop()

	for i := 0; i < 100; i++ {
		identifier := fmt.Sprintf("user-%d", i)
		rt.IsRequestAllowed(identifier)
	}
	
	time.Sleep(time.Millisecond * 300)
	
	expectedSize := 0
	mapSize := rt.GetStorageSize()

	if expectedSize != mapSize {
		t.Errorf("storage of the rate limiter must be cleaned up after window expires, expected %d, got %d", expectedSize, mapSize)
	}
}

func TestCleanupDoesNotAffectActiveRequests(t *testing.T) {
	rt := NewRateLimiter(1, time.Millisecond* 100)
	defer rt.Stop()

	for i := 0; i < 10; i++ {
		rt.IsRequestAllowed(fmt.Sprintf("old-user-%d", i))
	}

	time.Sleep(time.Millisecond * 300)
	rt.cleanup()

	rt.IsRequestAllowed("user-ege")
	rt.IsRequestAllowed("user-ece")

	expected := 2
	got := rt.GetStorageSize()

	if got != expected {
		t.Errorf("after cleaning up there must be 2 entries left, got %d expected %d", got, expected)
	}
}