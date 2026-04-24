package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"samll-trading-back/api/config"
	"samll-trading-back/api/database"
	"samll-trading-back/api/domains"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

var jwks keyfunc.Keyfunc

func InitJWKS() error {
	jwksURL := fmt.Sprintf("%s/auth/v1/.well-known/jwks.json", config.GetSupabaseURL())

	var err error
	jwks, err = keyfunc.NewDefault([]string{jwksURL})
	if err != nil {
		return fmt.Errorf("inicializar JWKS desde %s: %w", jwksURL, err)
	}

	log.Println("✅ Sistema de validación JWKS (Asimétrico) inicializado")
	return nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		if jwks == nil {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"error": "Servicio de autenticación no disponible"})
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, jwks.Keyfunc)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token signature"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		userID, ok := claims["sub"].(string)
		if !ok || userID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
			return
		}

		var user domains.User
		db := database.GetDB()

		if err := db.First(&user, "id = ?", userID).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				log.Printf("AuthMiddleware DB error userID=%s: %v", userID, err)
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "UNAUTHORIZED"})
			return
		}

		if !user.IsActive {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "ACCOUNT_PENDING",
				"message": "Tu cuenta está pendiente de aprobación por el administrador.",
			})
			return
		}

		c.Set("currentUser", user)
		c.Set("userID", user.ID)
		c.Next()
	}
}
