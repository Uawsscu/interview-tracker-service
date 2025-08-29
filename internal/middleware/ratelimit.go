package middleware

import (
	"net"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type limiterKey struct {
	IP   string
	Path string
}

type store struct {
	mu       sync.Mutex
	limiters map[limiterKey]*entry
	ttl      time.Duration
}

type entry struct {
	limiter *rate.Limiter
	last    time.Time
}

func newStore(ttl time.Duration) *store {
	s := &store{
		limiters: make(map[limiterKey]*entry),
		ttl:      ttl,
	}
	// background cleanup
	go func() {
		ticker := time.NewTicker(ttl)
		defer ticker.Stop()
		for range ticker.C {
			s.gc()
		}
	}()
	return s
}

func (s *store) get(ip, path string, r rate.Limit, burst int) *rate.Limiter {
	key := limiterKey{IP: ip, Path: path}
	now := time.Now()

	s.mu.Lock()
	defer s.mu.Unlock()

	if e, ok := s.limiters[key]; ok {
		e.last = now
		return e.limiter
	}
	lim := rate.NewLimiter(r, burst)
	s.limiters[key] = &entry{limiter: lim, last: now}
	return lim
}

func (s *store) gc() {
	cutoff := time.Now().Add(-s.ttl)
	s.mu.Lock()
	for k, e := range s.limiters {
		if e.last.Before(cutoff) {
			delete(s.limiters, k)
		}
	}
	s.mu.Unlock()
}

func clientIP(c *gin.Context) string {
	ip := c.ClientIP()
	if ip != "" {
		return ip
	}

	h := c.Request.Header.Get("X-Forwarded-For")
	if h != "" {
		parts := strings.Split(h, ",")
		return strings.TrimSpace(parts[0])
	}
	host, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err == nil && host != "" {
		return host
	}
	return c.Request.RemoteAddr
}

// RateLimitPerMinute จำกัดจำนวนการเรียก API ต่อ IP ต่อเส้นทาง
//   - limit: จำนวนครั้งสูงสุดต่อ 1 นาที (เช่น 60)
//   - burst: จำนวนครั้งเกินพิเศษที่อนุญาต (เช่น 20)
//   - ttl: เวลาที่เก็บสถานะ limiter ไว้ในหน่วยความจำ (เช่น 10 นาที)
func RateLimitPerMinute(limit int, burst int, ttl time.Duration) gin.HandlerFunc {
	if limit <= 0 {
		limit = 60
	}
	if burst < 1 {
		burst = 1
	}
	if ttl <= 0 {
		ttl = 10 * time.Minute
	}
	ratePerToken := rate.Every(time.Minute / time.Duration(limit))
	st := newStore(ttl)

	return func(c *gin.Context) {
		ip := clientIP(c)
		path := c.FullPath()
		if path == "" {
			// for routes added without named path (shouldn’t happen normally)
			path = c.Request.URL.Path
		}
		lim := st.get(ip, path, ratePerToken, burst)
		if !lim.Allow() {
			c.AbortWithStatusJSON(429, gin.H{
				"error":       "rate_limited",
				"message":     "Too many requests. Please slow down.",
				"retry_after": 1, // seconds (hint only)
			})
			return
		}
		c.Next()
	}
}
