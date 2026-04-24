package accounts

import (
	"log"
	"net/http"

	"samll-trading-back/api/database"
	"samll-trading-back/api/response"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type accountListItem struct {
	ID             string          `gorm:"column:id"              json:"id"`
	Name           string          `gorm:"column:name"            json:"name"`
	Broker         string          `gorm:"column:broker"          json:"broker"`
	AccountType    string          `gorm:"column:account_type"    json:"account_type"`
	Status         string          `gorm:"column:status"          json:"status"`
	Currency       string          `gorm:"column:currency"        json:"currency"`
	InitialBalance decimal.Decimal `gorm:"column:initial_balance" json:"initial_balance"`
	CurrentBalance decimal.Decimal `gorm:"column:current_balance" json:"current_balance"`
	TotalPnL       decimal.Decimal `gorm:"column:total_pnl"       json:"total_pnl"`
}

func ListAccounts(c *gin.Context) {
	userID := c.GetString("userID")

	var items []accountListItem
	db := database.GetDB()

	err := db.Raw(`
		SELECT
			ta.id,
			ta.name,
			ta.broker,
			ta.account_type,
			ta.status,
			ta.currency,
			ta.initial_balance,
			ta.initial_balance
				+ COALESCE(SUM(t.pnl),        0::numeric)
				- COALESCE(SUM(t.commission),  0::numeric) AS current_balance,
			COALESCE(SUM(t.pnl),        0::numeric)
				- COALESCE(SUM(t.commission),  0::numeric) AS total_pnl
		FROM trading_accounts ta
		LEFT JOIN trades t
			ON t.account_id = ta.id
			AND t.status    = 'CLOSED'
			AND t.deleted_at IS NULL
		WHERE ta.user_id    = ?
		  AND ta.deleted_at IS NULL
		GROUP BY ta.id
		ORDER BY ta.created_at DESC
	`, userID).Scan(&items).Error

	if err != nil {
		log.Printf("ListAccounts db error userID=%s: %v", userID, err)
		response.InternalError(c)
		return
	}

	if items == nil {
		items = []accountListItem{}
	}

	c.JSON(http.StatusOK, gin.H{"accounts": items})
}
