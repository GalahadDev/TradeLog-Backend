package admin

import (
	"errors"
	"log"
	"net/http"

	"samll-trading-back/api/database"
	"samll-trading-back/api/domains"
	"samll-trading-back/api/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user domains.User

	db := database.GetDB()

	if err := db.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "Usuario no encontrado")
		} else {
			log.Printf("AdminGetUserByID db error userID=%s: %v", id, err)
			response.InternalError(c)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
