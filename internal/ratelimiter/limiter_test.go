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
}


