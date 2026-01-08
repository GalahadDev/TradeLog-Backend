package trades

import (
	"net/http"
	"samll-trading-back/api/database"
	"samll-trading-back/api/domains"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetTrades(c *gin.Context) {
	userID := c.GetString("userID")

	// Paginación
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	var trades []domains.Trade
	var total int64

	db := database.GetDB()

	// Query Base con Tenant Isolation
	query := db.Model(&domains.Trade{}).Where("user_id = ?", userID)

	// Contar total
	query.Count(&total)

	// Obtener datos paginados y ordenados por fecha entrada
	if err := query.Order("entry_date desc").Offset(offset).Limit(limit).Find(&trades).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener trades"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": trades,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

func GetTradeByID(c *gin.Context) {
	// 1. Obtener IDs
	userID := c.GetString("userID")
	tradeID := c.Param("id")

	var trade domains.Trade
	db := database.GetDB()

	// 2. Buscar con Tenant Isolation
	// Solo devolvemos el trade si el ID coincide Y el dueño es el usuario actual
	if err := db.Where("id = ? AND user_id = ?", tradeID, userID).First(&trade).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trade no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"trade": trade})
}
