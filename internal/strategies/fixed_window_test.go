package strategies


import (
	"time"
	"testing"
)
func TestIsRequestAllowed(t *testing.T) {
	t.Run("single request", func(t *testing.T) {
		rt := &FixedWindowStrategy{limit: 1, storage: map[string]WindowData{}, windowSize: time.Minute}
		allowed, _ := rt.IsRequestAllowed("user123")

		if !allowed {
			t.Error("First request should be allowed")
		}

	})

	t.Run("2 requests at once", func(t *testing.T) {
		rt := &FixedWindowStrategy{limit: 1, storage: map[string]WindowData{}, windowSize: time.Minute}
		rt.IsRequestAllowed("user123")
		allowed, _ := rt.IsRequestAllowed("user123")

		if allowed {
			t.Error("Second request should be disallowed when limit is 1")
		}
	})
}

func TestFixedWindowStrategy(t *testing.T) {
	t.Run("single request", func(t *testing.T) {
		rt := &FixedWindowStrategy{limit: 1, storage: map[string]WindowData{}, windowSize: time.Minute}
		allowed, _ := rt.IsRequestAllowed("user123")
	
		if !allowed {
			t.Error("First request should be allowed")
		}
	
	})
	
	t.Run("2 requests at once", func(t *testing.T) {
		rt := &FixedWindowStrategy{limit: 1, storage: map[string]WindowData{}, windowSize: time.Minute}
		rt.IsRequestAllowed("user123")
		allowed, _ := rt.IsRequestAllowed("user123")
	
		if allowed {
			t.Error("Second request should be disallowed when limit is 1")
		}
	})

}
