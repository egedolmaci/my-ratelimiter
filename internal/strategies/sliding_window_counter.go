package strategies

import (
	"sync"
	"time"
)

type SlidingWindowCounterStrategy struct {
	limit int
	windowSize time.Duration
	timeProvider TimeProvider
	storage map[string]Data
	mu sync.RWMutex
	cleanupInterval time.Duration
	stopCleanup chan struct{}
	cleanupDone chan struct{}
}

type Data struct {
	currentWindow WindowData
	prevWindow WindowData
}

func NewSlidingWindowCountStrategy(limit int, windowSize time.Duration, timeProvider TimeProvider ) *SlidingWindowCounterStrategy {
	f := &SlidingWindowCounterStrategy{
		limit:limit,
		windowSize: windowSize,
		timeProvider: timeProvider,
		storage: make(map[string]Data),
		cleanupInterval: 185 * time.Second,
		cleanupDone: make(chan struct{}),
		stopCleanup: make(chan struct{}),
	}

	go f.startCleanup()
	return f
}

func (s *SlidingWindowCounterStrategy) IsRequestAllowed(identifier string) (bool, int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, exists := s.storage[identifier]
	if !exists {
		s.storage[identifier] = Data{currentWindow: WindowData{count: 1, timestamp: s.timeProvider.Now()}} 
		return true, s.limit - 1
	}

	currentWindowStart := s.timeProvider.Now().Truncate(s.windowSize)
	if currentWindowStart.After(data.currentWindow.timestamp) {
		data.prevWindow = data.currentWindow
		data.currentWindow = WindowData{count: 1, timestamp: currentWindowStart}
		s.storage[identifier] = data
		return true, s.limit - 1
	}

	timeElapsed := s.timeProvider.Now().Sub(currentWindowStart)
	percentageElapsed := float64(timeElapsed) / float64(s.windowSize)
	weight := 1.0 - percentageElapsed
	weightedLimit := float64(data.prevWindow.count) * weight + float64(data.currentWindow.count)

	if weightedLimit >= float64(s.limit) {
		return false, 0
	} 


	data.currentWindow.count++
	s.storage[identifier] = data
	remaining := s.limit - int(weightedLimit) - 1
	return true, remaining

}

func (s *SlidingWindowCounterStrategy) Stop() {
	close(s.stopCleanup)
	<-s.cleanupDone
}

func (s *SlidingWindowCounterStrategy) startCleanup() {
	ticker := time.NewTicker(s.cleanupInterval)

	for {
		select {
		case <-ticker.C:
			s.cleanup()
		case <-s.stopCleanup:
			close(s.cleanupDone)
			return
		}
	}
}

func (s * SlidingWindowCounterStrategy) cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	currentWindowStart := s.timeProvider.Now().Truncate(s.windowSize)
	for identifier, data := range s.storage {
		if currentWindowStart.After(data.currentWindow.timestamp) {
			delete(s.storage, identifier)
		}
	}
}