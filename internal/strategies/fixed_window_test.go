package strategies

import (
	"fmt"
	"testing"
	"time"
)

type MockTimeProvider struct {
	currentTime time.Time
}

func (m *MockTimeProvider) Now() time.Time {
	return m.currentTime
}

func (m *MockTimeProvider) Advance(pass time.Duration) {
	m.currentTime = m.currentTime.Add(pass)
}

func TestFixedWindowStrategy(t *testing.T) {
	t.Run("single request", func(t *testing.T) {
  		strategy := NewFixedWindowStrategy(1, 100 * time.Millisecond, &MockTimeProvider{currentTime: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)})
		 allowed, _ := strategy.IsRequestAllowed("user123")
	
		if !allowed {
			t.Error("First request should be allowed")
		}
	
	})
	
	t.Run("2 requests at once", func(t *testing.T) {
  		strategy := NewFixedWindowStrategy(1, 100 * time.Millisecond, &MockTimeProvider{currentTime: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)})
		strategy.IsRequestAllowed("user123")
		allowed, _ := strategy.IsRequestAllowed("user123")
	
		if allowed {
			t.Error("First request should be allowed")
		}
	})
}	
			
func TestCleanup(t *testing.T) {
	mockTimeProvider := &MockTimeProvider{currentTime: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)}
    strategy := NewFixedWindowStrategy(1, 100 * time.Millisecond, mockTimeProvider) 
    defer strategy.Stop()

    strategy.IsRequestAllowed("user1")

    if strategy.getStorageSize() != 1 {
        t.Error("Should have 1 entry")
    }

	mockTimeProvider.Advance(250 * time.Millisecond)
	strategy.cleanup()

    if strategy.getStorageSize() != 0 {
        t.Error("Should be cleaned up")
    }
}

func TestCleanupDoesNotAffectActiveRequests(t *testing.T) {
	mockTimeProvider := &MockTimeProvider{currentTime: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)}
    strategy := NewFixedWindowStrategy(1, 100 * time.Millisecond, mockTimeProvider)

	for i := range 10 {
    	strategy.IsRequestAllowed(fmt.Sprintf("old-user-%d", i))
	}

	mockTimeProvider.Advance(250 * time.Millisecond)

	strategy.IsRequestAllowed("new-user-ege")

	strategy.cleanup()

    if strategy.getStorageSize() != 1 {
        t.Error("Should have 1 entry")
    }

}
