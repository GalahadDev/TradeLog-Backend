package admin

import (
	"net/http"
	"samll-trading-back/api/database"
	"samll-trading-back/api/domains"

	"github.com/gin-gonic/gin"
)

// Usamos punteros (*) para saber si el campo vino en el JSON o no.
// Si es nil, no se actualiza.
type AdminPatchReq struct {
	FullName          *string `json:"full_name"`
	PhoneNumber       *string `json:"phone_number"`
	Bio               *string `json:"bio"`
	TradingExperience *string `json:"trading_experience"`
	AvatarURL         *string `json:"avatar_url"`
	Role              *string `json:"role"`
	IsActive          *bool   `json:"is_active"`
	IsVerified        *bool   `json:"is_verified"`
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req AdminPatchReq

	// Validar JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: revisa los tipos de datos (ej: bool sin comillas)"})
		return
	}

	db := database.GetDB()
	var user domains.User

	if err := db.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	// Actualización Dinámica: Solo actualizamos lo que no sea nil
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

	// Updates de GORM actualiza solo los campos presentes en el mapa
	if err := db.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar cambios"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Usuario actualizado", "user": user})
}
