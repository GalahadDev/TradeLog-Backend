package middleware

import (
	"log/slog"
	"time"

	"samll-trading-back/api/logger"

	"github.com/gin-gonic/gin"
)

// RequestLogger registra cada request/response HTTP usando el logger estructurado slog.
// Debe ejecutarse después de RequestID para que request_id ya esté disponible en el contexto.
//
// Campos registrados: method, path, status, latency_ms, ip, request_id, user_agent.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		rid := c.GetString(RequestIDKey)

		level := slog.LevelInfo
		if status >= 500 {
			level = slog.LevelError
		} else if status >= 400 {
			level = slog.LevelWarn
		}

		logger.L().Log(
			c.Request.Context(),
			level,
			"http",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", status),
			slog.Int64("latency_ms", latency.Milliseconds()),
			slog.String("ip", c.ClientIP()),
			slog.String("request_id", rid),
			slog.String("user_agent", c.Request.UserAgent()),
		)
	}
}
