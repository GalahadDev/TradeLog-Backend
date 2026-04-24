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

type AdminPatchReq struct {
	FullName          *string           `json:"full_name"`
	PhoneNumber       *string           `json:"phone_number"`
	Bio               *string           `json:"bio"`
	TradingExperience *string           `json:"trading_experience"`
	AvatarURL         *string           `json:"avatar_url"`
	Role              *domains.UserRole `json:"role"`
	IsActive          *bool             `json:"is_active"`
	IsVerified        *bool             `json:"is_verified"`
}

var allowedRoles = map[domains.UserRole]bool{domains.RoleUser: true, domains.RoleAdmin: true}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req AdminPatchReq

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("AdminUpdateUser bind error: %v", err)
		response.BadRequest(c, "Datos inválidos: revisa los tipos de datos (ej: bool sin comillas)")
		return
	}

	if req.Role != nil && !allowedRoles[*req.Role] {
		response.BadRequest(c, "role debe ser 'user' o 'admin'")
		return
	}

	db := database.GetDB()
	var user domains.User

	if err := db.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "Usuario no encontrado")
		} else {
			log.Printf("AdminUpdateUser db error userID=%s: %v", id, err)
			response.InternalError(c)
		}
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
	if req.Role != nil {
		updates["role"] = *req.Role
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if req.IsVerified != nil {
		updates["is_verified"] = *req.IsVerified
	}

	if err := db.Model(&user).Updates(updates).Error; err != nil {
		log.Printf("AdminUpdateUser update error userID=%s: %v", id, err)
		response.InternalError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario actualizado", "user": user})
}
