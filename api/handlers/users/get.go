package users

import (
	"net/http"
	"samll-trading-back/api/domains"

	"github.com/gin-gonic/gin"
)

func GetMe(c *gin.Context) {

	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No user context found"})
		return
	}

	currentUser, ok := user.(domains.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error interno de autenticación"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":    currentUser,
		"message": "Autenticación exitosa. Bienvenido al Trading Journal.",
	})
}
