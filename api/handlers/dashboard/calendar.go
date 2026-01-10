package dashboard

import (
	"net/http"
	"samll-trading-back/api/database"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

// CalendarDay: Estructura optimizada para la respuesta JSON.

type CalendarDay struct {
	Date       string          `json:"date"`
	TotalPnL   decimal.Decimal `json:"total_pnl" gorm:"column:total_pnl"`
	TradeCount int             `json:"trade_count" gorm:"column:trade_count"`
}

func GetCalendarMetrics(c *gin.Context) {
	userID := c.GetString("userID")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	db := database.GetDB()
	var results []CalendarDay

	// CONSTRUCCIÓN DE LA CONSULTA OPTIMIZADA
	query := db.Table("trades").
		Select("TO_CHAR(entry_date, 'YYYY-MM-DD') as date, SUM(pnl - commission) as total_pnl, COUNT(*) as trade_count").
		Where("user_id = ?", userID).
		Where("status = ?", "CLOSED")

	// 2. Aplicamos filtros de fecha si existen
	if startDate != "" && endDate != "" {
		query = query.Where("entry_date >= ? AND entry_date <= ?", startDate, endDate)
	}

	// 3. Agrupación y Orden
	err := query.Group("TO_CHAR(entry_date, 'YYYY-MM-DD')").
		Order("date ASC").
		Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error calculando métricas de calendario", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": results,
	})
}
