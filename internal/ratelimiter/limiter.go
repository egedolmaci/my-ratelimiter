package ratelimiter

import (
	"sync"
	"time"
)

type Ratelimiter struct {
	limit      int
	storage    map[string]WindowData
	windowSize time.Duration
	mu sync.RWMutex

	stopCleanup chan struct{}
	cleanupDone chan struct{}
	cleanupInterval time.Duration
}

type WindowData struct {
	count     int
	timestamp time.Time
}

func (r *Ratelimiter) IsRequestAllowed(identifier string) (bool, int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	data, exists := r.storage[identifier]
	if !exists || time.Now().After(data.timestamp.Add(r.windowSize)) {
		r.storage[identifier] = WindowData{count: 1, timestamp: time.Now()}
		return true, r.limit - 1
	}

	if data.count < r.limit {
		data.count++
		r.storage[identifier] = data
		return true, r.limit - data.count
	}

	return false, 0
}

func NewRateLimiter(limit int, windowSize time.Duration) *Ratelimiter {
	rt := &Ratelimiter{
		limit:      limit,
		storage:    map[string]WindowData{},
		windowSize: windowSize,

		cleanupInterval: windowSize * 2,
		stopCleanup: make(chan struct{}),
		cleanupDone: make(chan struct{}),
	}

	go rt.startCleanup()

	return rt
}

func (r *Ratelimiter) GetStorageSize() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.storage)
}

func (r *Ratelimiter) cleanup() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for identifier, data := range r.storage {
		if time.Now().After(data.timestamp.Add(r.windowSize)) {
			delete(r.storage, identifier)
		}
	}
}

func (r *Ratelimiter) startCleanup() {
	ticker := time.NewTicker(r.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.cleanup()
		case <-r.stopCleanup:
			close(r.cleanupDone)
			return
		}
	}
}

func (r *Ratelimiter) Stop() {
	close(r.stopCleanup)
	<-r.cleanupDone
}