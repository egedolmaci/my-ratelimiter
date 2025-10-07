package strategies

import (
	"sync"
	"time"
)

type SlidingWindowLogStrategy struct {
	limit        int
	windowSize   time.Duration
	storage      map[string][]time.Time
	timeProvider TimeProvider
	mu           sync.RWMutex
	stopCleanup  chan struct{}
	cleanupDone  chan struct{}
}

type RealTimeProvider struct{}

func (r *RealTimeProvider) Now() time.Time {
	return time.Now()
}

func NewSlidingWindowLogStrategy(limit int, windowSize time.Duration, timeProvider TimeProvider) *SlidingWindowLogStrategy {
	strategy := &SlidingWindowLogStrategy{
		limit:        limit,
		windowSize:   windowSize,
		storage:      map[string][]time.Time{},
		timeProvider: timeProvider,
		stopCleanup:  make(chan struct{}),
		cleanupDone:  make(chan struct{}),
	}

	go strategy.startCleanup()

	return strategy
}

func (s *SlidingWindowLogStrategy) IsRequestAllowed(identifier string) (bool, int) {
	return s.checkStorage(identifier)
}

func (s *SlidingWindowLogStrategy) checkStorage(identifier string) (bool, int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, exists := s.storage[identifier]

	if !exists {
		data = append(data, s.timeProvider.Now())
		s.storage[identifier] = data
		return true, s.limit - 1
	}

	newList := s.cleanStorage(data)
	if len(newList) < s.limit {
		newList = append(newList, s.timeProvider.Now())
		s.storage[identifier] = newList
		return true, s.limit - len(newList)
	}

	return false, 0
}

func (s *SlidingWindowLogStrategy) cleanStorage(list []time.Time) []time.Time {
	count := 0
	for _, elm := range list {
		if s.timeProvider.Now().After(elm.Add(s.windowSize)) {
			count++
		}
	}

	return list[count:]
}

func (s *SlidingWindowLogStrategy) startCleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

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

func (s *SlidingWindowLogStrategy) cleanup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for identifier, list := range s.storage {
		newList := s.cleanStorage(list)
		if len(newList) == 0 {
			delete(s.storage, identifier)
		} else {
			s.storage[identifier] = newList
		}
	}
}

func (s *SlidingWindowLogStrategy) Stop() {
	close(s.stopCleanup)
	<-s.cleanupDone
}
