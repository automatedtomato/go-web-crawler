package ratelimiter

import (
	"context"
	"sync"

	"golang.org/x/time/rate"
)

// Host別レート制限を管理
type HostLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	r        rate.Limit
	b        int
}

func NewHostLimiter(r rate.Limit, b int) *HostLimiter {
	return &HostLimiter{
		limiters: make(map[string]*rate.Limiter),
		r:        r, // Num of requests per second
		b:        b, // Burst size (num of request allowed in a instant period of time)
	}
}

// Retrieve rate limiter for host
func (l *HostLimiter) GetLimiter(host string) *rate.Limiter {
	l.mu.RLock()
	limiter, exist := l.limiters[host]
	defer l.mu.RUnlock()

	if !exist {
		l.mu.Lock()
		// Double-check if other goroutine already created
		limiter, exists := l.limiters[host]
		if !exists {
			limiter = rate.NewLimiter(l.r, l.b)
			l.limiters[host] = limiter
		}
		l.mu.Unlock()
	}

	return limiter
}

func (l *HostLimiter) Wait(host string) {
	limiter := l.GetLimiter(host)
	_ = limiter.Wait(context.Background())
}

func (l *HostLimiter) Allow(host string) bool {
	limiter := l.GetLimiter(host)
	return limiter.Allow()
}
