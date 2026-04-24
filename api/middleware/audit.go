package middleware

import (
	"log/slog"
	"net/http"

	"samll-trading-back/api/database"
	"samll-trading-back/api/domains"
	"samll-trading-back/api/logger"

	"github.com/gin-gonic/gin"
)

// AuditAdmin registra cada operación administrativa que modifica estado (no-GET)
// y que finaliza con código 2xx en la tabla audit_logs.
//
// Debe ubicarse después de AuthMiddleware para que "currentUser" esté disponible
// en el contexto de Gin.
//
// targetType identifica el recurso sobre el que se opera (ej. "user").
func AuditAdmin(targetType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Request.Method == http.MethodGet {
			return
		}
		status := c.Writer.Status()
		if status < http.StatusOK || status >= http.StatusMultipleChoices {
			return
		}

		user, exists := c.Get("currentUser")
		if !exists {
			return
		}
		currentUser, ok := user.(domains.User)
		if !ok {
			return
		}

		entry := domains.AuditLog{
			AdminID:    currentUser.ID,
			Action:     c.Request.Method + " " + c.FullPath(),
			TargetType: targetType,
			TargetID:   c.Param("id"),
			IP:         c.ClientIP(),
			RequestID:  c.GetString(RequestIDKey),
		}

		if err := database.GetDB().Create(&entry).Error; err != nil {
			logger.L().Error("audit log write failed",
				slog.String("action", entry.Action),
				slog.String("admin_id", entry.AdminID),
				slog.String("request_id", entry.RequestID),
				slog.String("err", err.Error()),
			)
		}
	}
}
