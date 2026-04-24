package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	limiters   = make(map[string]*ipLimiter)
	limitersMu sync.Mutex
)

func init() {
	// Goroutine de limpieza: elimina entradas inactivas cada 5 minutos
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			limitersMu.Lock()
			for ip, l := range limiters {
				if time.Since(l.lastSeen) > 10*time.Minute {
					delete(limiters, ip)
				}
			}
			limitersMu.Unlock()
		}
	}()
}

func getLimiter(ip string) *rate.Limiter {
	limitersMu.Lock()
	defer limitersMu.Unlock()

	if l, exists := limiters[ip]; exists {
		l.lastSeen = time.Now()
		return l.limiter
	}

	// 60 peticiones por minuto con burst de 20
	l := &ipLimiter{
		limiter:  rate.NewLimiter(rate.Every(time.Second), 20),
		lastSeen: time.Now(),
	}
	limiters[ip] = l
	return l.limiter
}

// RateLimit aplica un límite de 60 req/min por IP (burst 20).
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !getLimiter(ip).Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":   "RATE_LIMIT_EXCEEDED",
				"message": "Demasiadas peticiones. Intenta de nuevo en un momento.",
			})
			return
		}
		c.Next()
	}
}
