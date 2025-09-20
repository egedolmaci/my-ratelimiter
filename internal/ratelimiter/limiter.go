package ratelimiter

import "time"

type Ratelimiter struct {
	limit   int
	storage map[string]WindowData
	windowSize time.Duration
}

type WindowData struct {
	count int
	timestamp time.Time
}

func (r *Ratelimiter) IsRequestAllowed(identifier string) (bool, int) {
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
	return &Ratelimiter{
		limit: limit,
		storage: map[string]WindowData{},
		windowSize: windowSize,
	}
}


