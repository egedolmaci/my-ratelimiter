package ratelimiter

import (
	"testing"
	"time"
)


func TestConcurrentAccess(t *testing.T) {

	limit := 20
	iters := 10000
	identifier := "127.0.0.1:1000"

	ratelimiter := NewRateLimiter(limit, time.Millisecond * 1)

	allowedChan := make(chan bool)
	startChan := make(chan struct{})
	isAllowedCount := 0; isNotAllowedCount := 0
	
	for range iters {
		go func(identifier string) {
			<- startChan
			isAllowed, _ := ratelimiter.IsRequestAllowed(identifier)
			allowedChan <- isAllowed
		}(identifier)
	}

	close(startChan)
	
	for range iters {
		isAllowedResult := <- allowedChan
		if isAllowedResult {
			isAllowedCount++
		} else {
			isNotAllowedCount++
		}
	}

	expectedAllowedCount := limit
	expectedNotAllowedCount := iters - limit
	if isAllowedCount != expectedAllowedCount || isNotAllowedCount != expectedNotAllowedCount {
		t.Errorf("expected allowed count %d got %d expected not allowed count %d got %d",
			expectedAllowedCount,
			isAllowedCount,
			expectedNotAllowedCount,
			isNotAllowedCount)
	}
}



