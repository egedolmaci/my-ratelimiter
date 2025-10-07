package strategies

import (
	"sync"
	"time"
)

type FixedWindowStrategy struct {
	limit      int
	storage    map[string]WindowData
	windowSize time.Duration
	mu         sync.RWMutex

	stopCleanup     chan struct{}
	cleanupDone     chan struct{}
	cleanupInterval time.Duration
	timeProvider    TimeProvider
}

type WindowData struct {
	count     int
	timestamp time.Time
}

func (f *FixedWindowStrategy) IsRequestAllowed(identifier string) (bool, int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	data, exists := f.storage[identifier]
	currentWindow := f.timeProvider.Now().Truncate(f.windowSize)
	if !exists || currentWindow != data.timestamp {
		f.storage[identifier] = WindowData{count: 1, timestamp: currentWindow}
		return true, f.limit - 1
	}

	if data.count < f.limit {
		data.count++
		f.storage[identifier] = data
		return true, f.limit - data.count
	}

	return false, 0
}

func (f *FixedWindowStrategy) Stop() {
	close(f.stopCleanup)
	<-f.cleanupDone
}

func NewFixedWindowStrategy(limit int, windowSize time.Duration, TimeProvider TimeProvider) *FixedWindowStrategy {
	f := &FixedWindowStrategy{
		limit:           limit,
		storage:         map[string]WindowData{},
		windowSize:      windowSize,
		cleanupInterval: windowSize * 2,
		stopCleanup:     make(chan struct{}),
		cleanupDone:     make(chan struct{}),
		timeProvider:    TimeProvider,
	}

	go f.startCleanup()
	return f
}

func (f *FixedWindowStrategy) cleanup() {
	f.mu.Lock()
	defer f.mu.Unlock()

	for identifier, data := range f.storage {
		currentWindow := f.timeProvider.Now().Truncate(f.windowSize)
		if currentWindow != data.timestamp {
			delete(f.storage, identifier)
		}
	}
}

func (f *FixedWindowStrategy) startCleanup() {
	ticker := time.NewTicker(f.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			f.cleanup()
		case <-f.stopCleanup:
			close(f.cleanupDone)
			return
		}
	}
}

func (f *FixedWindowStrategy) getStorageSize() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(f.storage)
}
