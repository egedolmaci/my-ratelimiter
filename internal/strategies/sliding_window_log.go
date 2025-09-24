package strategies

import "time"

type SlidingWindowLogStrategy struct {
	limit int
	windowSize time.Duration
}

func NewSlidingWindowLogStrategy(limit int, windowSize time.Duration) *SlidingWindowLogStrategy {{
	return &SlidingWindowLogStrategy{limit, windowSize}
}}

func (s *SlidingWindowLogStrategy) IsRequestAllowed(identifier string) (bool, int) {
	return true, 0
}

func (s *SlidingWindowLogStrategy) Stop() {}