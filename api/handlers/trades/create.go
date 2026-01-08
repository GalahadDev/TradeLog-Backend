package trades

import (
	"net/http"
	"strings"

	"samll-trading-back/api/database"
	"samll-trading-back/api/domains"

	"github.com/gin-gonic/gin"
)

func CreateTrade(c *gin.Context) {
	// 1. Tenant Isolation: Obtener ID del usuario autenticado
	userID := c.GetString("userID")

	var trade domains.Trade
	if err := c.ShouldBindJSON(&trade); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
		return
	}

	// 2. Sobrescribir el UserID para evitar suplantación
	trade.UserID = userID

	// 3. NORMALIZACIÓN
	trade.Direction = strings.ToUpper(trade.Direction)
	trade.Status = strings.ToUpper(trade.Status)

	// 4. Guardar en DB
	db := database.GetDB()
	if err := db.Create(&trade).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar el trade", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"trade": trade})
}
