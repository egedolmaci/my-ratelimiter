package ratelimiter

import (
	"time"

	"github.com/egedolmaci/my-ratelimiter/internal/strategies"
)

type RateLimitStrategy interface {
	IsRequestAllowed(identifier string) (bool, int)
	Stop()
}

type Ratelimiter struct {
	strategy RateLimitStrategy
}

func NewRateLimiterWithStrategy(strategy RateLimitStrategy) *Ratelimiter {
	return &Ratelimiter{strategy: strategy}
}

func (r *Ratelimiter) IsRequestAllowed(identifier string) (bool, int) {
	return r.strategy.IsRequestAllowed(identifier)
}

func (r *Ratelimiter) Stop() {
	r.strategy.Stop()
}

func NewRateLimiter(limit int, windowSize time.Duration) *Ratelimiter {
	strategy := strategies.NewFixedWindowStrategy(limit, windowSize)
	return NewRateLimiterWithStrategy(strategy)
}

