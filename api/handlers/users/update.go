package users

import (
	"errors"
	"log"
	"net/http"
	"net/url"

	"samll-trading-back/api/database"
	"samll-trading-back/api/domains"
	"samll-trading-back/api/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		log.Printf("UpdateMyProfile bind error: %v", err)
		response.BadRequest(c, "Datos inválidos")
		return
	}

	if req.FullName != nil && len(*req.FullName) > 100 {
		response.BadRequest(c, "full_name no puede superar 100 caracteres")
		return
	}
	if req.Bio != nil && len(*req.Bio) > 500 {
		response.BadRequest(c, "bio no puede superar 500 caracteres")
		return
	}
	if req.PhoneNumber != nil && len(*req.PhoneNumber) > 30 {
		response.BadRequest(c, "phone_number no puede superar 30 caracteres")
		return
	}
	if req.AvatarURL != nil && *req.AvatarURL != "" {
		parsed, err := url.Parse(*req.AvatarURL)
		if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
			response.BadRequest(c, "avatar_url debe ser una URL válida con esquema http o https")
			return
		}
	}

	db := database.GetDB()
	var user domains.User

	if err := db.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "Usuario no encontrado")
		} else {
			log.Printf("UpdateMyProfile db error userID=%s: %v", userID, err)
			response.InternalError(c)
		}
		return
	}

	updates := make(map[string]any)

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
		log.Printf("UpdateMyProfile update error userID=%s: %v", userID, err)
		response.InternalError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Perfil actualizado correctamente", "user": user})
}
