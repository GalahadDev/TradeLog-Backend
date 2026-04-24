package dashboard

import (
	"log"
	"net/http"
	"time"

	"samll-trading-back/api/database"
	"samll-trading-back/api/response"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type CalendarDay struct {
	Date       string          `json:"date"`
	TotalPnL   decimal.Decimal `json:"total_pnl"    gorm:"column:total_pnl"`
	TradeCount int             `json:"trade_count"  gorm:"column:trade_count"`
}

func GetCalendarMetrics(c *gin.Context) {
	accountID := c.GetString("accountID")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate != "" {
		if _, err := time.Parse("2006-01-02", startDate); err != nil {
			response.BadRequest(c, "Formato de start_date inválido. Usar YYYY-MM-DD")
			return
		}
	}
	if endDate != "" {
		if _, err := time.Parse("2006-01-02", endDate); err != nil {
			response.BadRequest(c, "Formato de end_date inválido. Usar YYYY-MM-DD")
			return
		}
	}

	db := database.GetDB()
	var results []CalendarDay

	query := db.Table("trades").
		Select("TO_CHAR(entry_date, 'YYYY-MM-DD') as date, SUM(pnl - commission) as total_pnl, COUNT(*) as trade_count").
		Where("account_id = ?", accountID).
		Where("status = ?", "CLOSED")

	if startDate != "" && endDate != "" {
		query = query.Where("entry_date >= ? AND entry_date <= ?", startDate, endDate)
	}

	err := query.Group("TO_CHAR(entry_date, 'YYYY-MM-DD')").
		Order("date ASC").
		Scan(&results).Error

	if err != nil {
		log.Printf("GetCalendarMetrics db error accountID=%s: %v", accountID, err)
		response.InternalError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}
