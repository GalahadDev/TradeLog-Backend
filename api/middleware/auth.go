package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"samll-trading-back/api/config"
	"samll-trading-back/api/database"
	"samll-trading-back/api/domains"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	jwks     keyfunc.Keyfunc
	jwksOnce sync.Once
)

// initJWKS inicializa el conjunto de claves.
func initJWKS() {
	jwksURL := fmt.Sprintf("%s/auth/v1/.well-known/jwks.json", config.GetSupabaseURL())

	var err error
	// Crea el JWKS. Automáticamente refresca las claves si cambian.
	jwks, err = keyfunc.NewDefault([]string{jwksURL})
	if err != nil {
		log.Fatalf("❌ Fallo al iniciar JWKS: %v", err)
	}
	log.Println("✅ Sistema de validación JWKS (Asimétrico) inicializado")
}

func AuthMiddleware() gin.HandlerFunc {
	// Aseguramos que el JWKS esté listo antes de aceptar peticiones
	jwksOnce.Do(initJWKS)

	return func(c *gin.Context) {
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

		// Extraer UUID
		userID, ok := claims["sub"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
			return
		}

		var user domains.User
		db := database.GetDB()

		// Buscamos al usuario
		if err := db.First(&user, "id = ?", userID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User profile not found in DB"})
			return
		}

		// Whitelist Check
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
