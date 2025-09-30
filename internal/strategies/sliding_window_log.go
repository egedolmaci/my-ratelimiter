package strategies

import "time"

type SlidingWindowLogStrategy struct {
	limit int
	windowSize time.Duration
	storage map[string]int
	timeProvider TimeProvider
}

type RealTimeProvider struct {}

func (r *RealTimeProvider) Now() time.Time {
	return time.Now()
}

func NewSlidingWindowLogStrategy(limit int, windowSize time.Duration, timeProvider TimeProvider) *SlidingWindowLogStrategy {{
	return &SlidingWindowLogStrategy{
		limit,
		windowSize,
		map[string]int{},
		&RealTimeProvider{},
	}
}}

func (s *SlidingWindowLogStrategy) IsRequestAllowed(identifier string) (bool, int) {
	return s.checkStorage(identifier)
}

func (s *SlidingWindowLogStrategy) checkStorage(identifier string) (bool, int) {
	for id, count := range s.storage {
		if id == identifier {
			if count < s.limit {
				s.addToStorage(identifier)
				return true, s.limit - count
			} else {
				return false, 0
			}
		} 
	}

	s.addToStorage(identifier)
	return true, s.limit 
}

func (s *SlidingWindowLogStrategy) addToStorage(identifier string) {
	s.storage[identifier]++
}

func (s *SlidingWindowLogStrategy) Stop() {}

