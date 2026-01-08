package middleware

import (
	"net/http"
	"samll-trading-back/api/domains"

	"github.com/gin-gonic/gin"
)

// AdminOnly verifica que el usuario en el contexto tenga rol 'admin'
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		userCtx, exists := c.Get("currentUser")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		user := userCtx.(domains.User)

		if user.Role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "INSUFFICIENT_PERMISSIONS",
				"message": "Acceso denegado: Se requieren privilegios de Administrador.",
			})
			return
		}
		c.Next()
	}
}
