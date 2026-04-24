package trades

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"samll-trading-back/api/database"
	"samll-trading-back/api/domains"
	"samll-trading-back/api/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetTrades(c *gin.Context) {
	accountID := c.GetString("accountID")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 1
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	var trades []domains.Trade
	var total int64

	db := database.GetDB()
	query := db.Model(&domains.Trade{}).Where("account_id = ?", accountID)
	query.Count(&total)

	if err := query.Order("entry_date desc").Offset(offset).Limit(limit).Find(&trades).Error; err != nil {
		log.Printf("GetTrades db error accountID=%s: %v", accountID, err)
		response.InternalError(c)
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
	accountID := c.GetString("accountID")
	tradeID := c.Param("trade_id")

	var trade domains.Trade
	db := database.GetDB()

	if err := db.Where("id = ? AND account_id = ?", tradeID, accountID).First(&trade).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "Trade no encontrado")
		} else {
			log.Printf("GetTradeByID db error tradeID=%s accountID=%s: %v", tradeID, accountID, err)
			response.InternalError(c)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"trade": trade})
}
