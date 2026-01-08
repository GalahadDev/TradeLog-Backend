package admin

import (
	"net/http"
	"samll-trading-back/api/database"
	"samll-trading-back/api/domains"

	"github.com/gin-gonic/gin"
)

func GetAllUsers(c *gin.Context) {
	db := database.GetDB()
	var users []domains.User

	// Filtros opcionales por Query Param
	// Ejemplo: /api/admin/users?pending=true
	pending := c.Query("pending")

	query := db.Model(&domains.User{})

	if pending == "true" {
		query = query.Where("is_active = ?", false)
	}

	// Ordenar por fecha de creación
	if err := query.Order("created_at desc").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}
