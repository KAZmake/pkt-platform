package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/KAZmake/pkt-platform/apps/api/assistant-svc/pkg/response"
	"golang.org/x/time/rate"
)

type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiter is a per-IP request rate limiter middleware.
type RateLimiter struct {
	mu       sync.Mutex
	limiters map[string]*ipLimiter
	rpm      int
}

// NewRateLimiter creates a middleware that allows at most rpm requests per minute per IP.
func NewRateLimiter(rpm int) *RateLimiter {
	rl := &RateLimiter{
		limiters: make(map[string]*ipLimiter),
		rpm:      rpm,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	entry, ok := rl.limiters[ip]
	if !ok {
		// rpm requests per minute → tokens per second = rpm/60, burst = rpm
		r := rate.Limit(float64(rl.rpm) / 60.0)
		entry = &ipLimiter{limiter: rate.NewLimiter(r, rl.rpm)}
		rl.limiters[ip] = entry
	}
	entry.lastSeen = time.Now()
	return entry.limiter
}

// Middleware returns an http.Handler middleware function.
func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}
		if !rl.getLimiter(ip).Allow() {
			response.TooManyRequests(w)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// cleanup removes entries that haven't been seen in the last 10 minutes.
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		for ip, entry := range rl.limiters {
			if time.Since(entry.lastSeen) > 10*time.Minute {
				delete(rl.limiters, ip)
			}
		}
		rl.mu.Unlock()
	}
}
