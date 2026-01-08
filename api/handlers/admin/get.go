package admin

import (
	"net/http"
	"samll-trading-back/api/database"
	"samll-trading-back/api/domains"

	"github.com/gin-gonic/gin"
)

func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user domains.User

	db := database.GetDB()
	// Buscamos por ID
	if err := db.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
