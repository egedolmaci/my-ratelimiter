package ratelimiter

import "time"

type Ratelimiter struct {
	limit   int
	storage map[string]WindowData
	windowSize time.Duration
}

type WindowData struct {
	count int
}

func (r *Ratelimiter) IsRequestAllowed(identifier string) bool {
	data, exists := r.storage[identifier]
	if !exists {
		r.storage[identifier] = WindowData{count: 1}
		return true
	}

	if data.count < r.limit {
		data.count++
		r.storage[identifier] = data
		return true
	}

	return false
}

func NewRateLimiter(limit int, windowSize time.Duration) *Ratelimiter {
	return &Ratelimiter{
		limit: limit,
		storage: map[string]WindowData{},
		windowSize: windowSize,
	}
}
