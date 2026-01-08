package users

import (
	"net/http"
	"samll-trading-back/api/domains"

	"github.com/gin-gonic/gin"
)

// GetMe devuelve la información del usuario actual extraída del contexto
func GetMe(c *gin.Context) {
	// 1. Recuperar el usuario que el Middleware guardó en el contexto
	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No user context found"})
		return
	}

	// 2. Castear la variable al tipo User
	currentUser := user.(domains.User)

	// 3. Responder con los datos
	c.JSON(http.StatusOK, gin.H{
		"user":    currentUser,
		"message": "Autenticación exitosa. Bienvenido al Trading Journal.",
	})
}
