package dashboard

import (
	"net/http"
	"samll-trading-back/api/database"
	"samll-trading-back/api/domains"
	"samll-trading-back/api/services/analytics"

	"github.com/gin-gonic/gin"
)

func GetStats(c *gin.Context) {
	userID := c.GetString("userID")

	// Filtros opcionales de fecha
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	db := database.GetDB()
	var trades []domains.Trade

	query := db.Where("user_id = ? AND status = ?", userID, "CLOSED").Order("exit_date asc")

	if startDate != "" && endDate != "" {
		query = query.Where("exit_date >= ? AND exit_date <= ?", startDate, endDate)
	}

	if err := query.Find(&trades).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo trades", "details": err.Error()})
		return
	}

	// 2. Llamar al servicio de cálculo
	stats := analytics.CalculateStats(trades)

	// 3. Responder
	c.JSON(http.StatusOK, gin.H{
		"stats": stats,
	})
}
