package dashboard

import (
	"log"
	"net/http"
	"time"

	"samll-trading-back/api/cache"
	"samll-trading-back/api/database"
	"samll-trading-back/api/domains"
	"samll-trading-back/api/response"
	"samll-trading-back/api/services/analytics"

	"github.com/gin-gonic/gin"
)

func GetStats(c *gin.Context) {
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

	cacheKey := cache.Key(accountID, startDate, endDate)
	if stats, ok := cache.Stats.Get(cacheKey); ok {
		c.JSON(http.StatusOK, gin.H{"stats": stats})
		return
	}

	db := database.GetDB()
	var trades []domains.Trade

	query := db.Where("account_id = ? AND status = ?", accountID, "CLOSED").Order("exit_date asc")

	if startDate != "" && endDate != "" {
		query = query.Where("exit_date >= ? AND exit_date <= ?", startDate, endDate)
	}

	if err := query.Find(&trades).Error; err != nil {
		log.Printf("GetStats db error accountID=%s: %v", accountID, err)
		response.InternalError(c)
		return
	}

	stats := analytics.CalculateStats(trades)
	cache.Stats.Set(cacheKey, stats)

	c.JSON(http.StatusOK, gin.H{"stats": stats})
}
