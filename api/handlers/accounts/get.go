package accounts

import (
	"log"
	"net/http"

	"samll-trading-back/api/database"
	"samll-trading-back/api/domains"
	"samll-trading-back/api/response"
	"samll-trading-back/api/services/analytics"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

func GetAccount(c *gin.Context) {
	accountID := c.GetString("accountID")
	db := database.GetDB()

	var account domains.TradingAccount
	if err := db.First(&account, "id = ?", accountID).Error; err != nil {
		log.Printf("GetAccount db error accountID=%s: %v", accountID, err)
		response.InternalError(c)
		return
	}

	var trades []domains.Trade
	if err := db.Where("account_id = ? AND status = 'CLOSED'", accountID).Find(&trades).Error; err != nil {
		log.Printf("GetAccount trades db error accountID=%s: %v", accountID, err)
		response.InternalError(c)
		return
	}

	stats := analytics.CalculateStats(trades)

	currentBalance := account.InitialBalance.Add(stats.TotalNetProfit)

	summary := gin.H{
		"current_balance": currentBalance,
		"total_pnl":       stats.TotalNetProfit,
		"total_trades":    stats.TotalTrades,
		"win_rate":        stats.WinRate,
	}

	if account.AccountType == "prop_firm" &&
		account.ProfitTarget != nil && !account.ProfitTarget.IsZero() &&
		account.MaxDrawdownLimit != nil {

		profitProgress := stats.TotalNetProfit.
			Div(*account.ProfitTarget).
			Mul(decimal.NewFromInt(100))

		drawdownUsed := stats.MaxDrawdown
		drawdownRemaining := account.MaxDrawdownLimit.Sub(drawdownUsed)

		summary["profit_progress"] = profitProgress
		summary["drawdown_used"] = drawdownUsed
		summary["drawdown_remaining"] = drawdownRemaining
	}

	c.JSON(http.StatusOK, gin.H{
		"account": account,
		"summary": summary,
	})
}
