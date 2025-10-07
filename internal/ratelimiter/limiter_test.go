package ratelimiter

import (
	"testing"
	"time"

	"github.com/egedolmaci/my-ratelimiter/internal/strategies"
)

type MockTimeProvider struct {
	currentTime time.Time
}

func (m *MockTimeProvider) Now() time.Time {
	return m.currentTime
}

func (m *MockTimeProvider) Advance(duration time.Duration) {
	m.currentTime = m.currentTime.Add(duration)
}

func TestIsRequestAllowedFixedWindow(t *testing.T) {
	t.Run("1 request fills limit and limit opens up for another", func(t *testing.T) {
		mockTimeProvider := &MockTimeProvider{}
		rt := NewRateLimiterWithStrategy(strategies.NewFixedWindowStrategy(1, time.Millisecond*100, mockTimeProvider))
		defer rt.Stop()
		rt.IsRequestAllowed("user123")

		mockTimeProvider.Advance(150 * time.Millisecond)

		allowed, _ := rt.IsRequestAllowed("user123")

		if !allowed {
			t.Errorf("It must be allowed since windowSize amount of time has passed")
		}
	})
	t.Run("1 request fills the limit and since enough time has not passed limit is full", func(t *testing.T) {
		mockTimeProvider := &MockTimeProvider{}
		rt := NewRateLimiterWithStrategy(strategies.NewFixedWindowStrategy(1, time.Millisecond*100, mockTimeProvider))
		defer rt.Stop()
		rt.IsRequestAllowed("user123")

		mockTimeProvider.Advance(10 * time.Millisecond)

		allowed, _ := rt.IsRequestAllowed("user123")

		if allowed {
			t.Errorf("It must not be allowed since required time for the limit to reset has not passed")
		}
	})

	t.Run("when a burst happens around the boundaries requests should be allowed", func(t *testing.T) {
		mockTimeProvider := &MockTimeProvider{}
		rl := NewRateLimiterWithStrategy(strategies.NewFixedWindowStrategy(10, time.Minute, mockTimeProvider))
		defer rl.Stop()

		mockTimeProvider.Advance(time.Second * 59)

		for range 10 {
			rl.IsRequestAllowed("ege")
		}

		mockTimeProvider.Advance(time.Second)

		for range 9 {
			rl.IsRequestAllowed("ege")
		}

		allowed, count := rl.IsRequestAllowed("ege")

		if !allowed {
			t.Errorf("request must be allowed expected %t %d, got %t, %d", true, 0, allowed, count)
		}
	})
}

func TestIsRequestAllowedSlidingWindowLog(t *testing.T) {
	t.Run("should deny request when all requests still in window", func(t *testing.T) {
		mockTimeProvider := &MockTimeProvider{}
		rl := NewRateLimiterWithStrategy(strategies.NewSlidingWindowLogStrategy(10, time.Minute, mockTimeProvider))
		defer rl.Stop()

		mockTimeProvider.Advance(30 * time.Second)

		for range 10 {
			rl.IsRequestAllowed("ege")
		}

		mockTimeProvider.Advance(30 * time.Second)

		allowed, _ := rl.IsRequestAllowed("ege")

		if allowed {
			t.Errorf("request should not be allowed expected false got %t", allowed)
		}
	})

	t.Run("should allow request after old requests expire from window", func(t *testing.T) {
		mockTimeProvider := &MockTimeProvider{}
		rl := NewRateLimiterWithStrategy(strategies.NewSlidingWindowLogStrategy(10, time.Minute, mockTimeProvider))
		defer rl.Stop()

		rl.IsRequestAllowed("ege")
		mockTimeProvider.Advance(55 * time.Second)

		for range 9 {
			rl.IsRequestAllowed("ege")
		}

		mockTimeProvider.Advance(10 * time.Second)

		allowed, remaining := rl.IsRequestAllowed("ege")

		if !allowed {
			t.Errorf("request should be allowed after window expires, got allowed=%t remaining=%d", allowed, remaining)
		}

		if remaining != 0 {
			t.Errorf("should have 0 remaining after first request in new window, got %d", remaining)
		}
	})
}

func TestRateLimiterWithConfig(t *testing.T) {
	t.Run("ratelimiter with config", func(t *testing.T) {
		config := Config{
			Strategy:   "fixed_window",
			Limit:      10,
			WindowSize: time.Minute,
		}

		rl := NewRatelimiterWithConfig(config)

		for i := 0; i < 10; i++ {
			allowed, _ := rl.IsRequestAllowed("ege")
			if !allowed {
				t.Errorf("%dth request should be allowed", i+1)
			}
		}

		allowed, _ := rl.IsRequestAllowed("ege")
		if allowed {
			t.Errorf("11th request should not be allowed")
		}

	})
}
