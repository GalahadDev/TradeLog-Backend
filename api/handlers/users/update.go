package users

import (
	"net/http"
	"samll-trading-back/api/database"
	"samll-trading-back/api/domains"

	"github.com/gin-gonic/gin"
)

type UserUpdateReq struct {
	FullName          *string `json:"full_name"`
	PhoneNumber       *string `json:"phone_number"`
	Bio               *string `json:"bio"`
	TradingExperience *string `json:"trading_experience"`
	AvatarURL         *string `json:"avatar_url"`
}

func UpdateMyProfile(c *gin.Context) {
	userID := c.GetString("userID")

	var req UserUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	db := database.GetDB()
	var user domains.User

	if err := db.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	updates := make(map[string]interface{})

	if req.FullName != nil {
		updates["full_name"] = *req.FullName
	}
	if req.PhoneNumber != nil {
		updates["phone_number"] = *req.PhoneNumber
	}
	if req.Bio != nil {
		updates["bio"] = *req.Bio
	}
	if req.TradingExperience != nil {
		updates["trading_experience"] = *req.TradingExperience
	}

	if req.AvatarURL != nil {
		updates["avatar_url"] = *req.AvatarURL
	}

	if err := db.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo actualizar el perfil"})
		return
	}

	// Devolvemos el usuario actualizado para que el front refresque la vista
	c.JSON(http.StatusOK, gin.H{
		"message": "Perfil actualizado correctamente",
		"user":    user,
	})
}
