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

type Config struct {
	Strategy   string
	Limit      int
	WindowSize time.Duration
}

func NewRatelimiterWithConfig(config Config) *Ratelimiter {
	return NewRateLimiter(config.Limit, config.WindowSize, &strategies.RealTimeProvider{}, config.Strategy)
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

func NewRateLimiter(limit int, windowSize time.Duration, timeProvider strategies.TimeProvider, strategyName string) *Ratelimiter {
	var strategy RateLimitStrategy
	if strategyName == "fixed_window" {
		strategy = strategies.NewFixedWindowStrategy(limit, windowSize, timeProvider)
	} else if strategyName == "sliding_window_log" {
		strategy = strategies.NewSlidingWindowLogStrategy(limit, windowSize, timeProvider)
	}

	return NewRateLimiterWithStrategy(strategy)
}
