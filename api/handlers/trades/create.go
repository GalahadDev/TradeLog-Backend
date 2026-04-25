package trades

import (
	"log"
	"net/http"
	"strings"

	"samll-trading-back/api/cache"
	"samll-trading-back/api/database"
	"samll-trading-back/api/domains"
	"samll-trading-back/api/response"

	"github.com/gin-gonic/gin"
)

var validDirections = map[string]bool{"LONG": true, "SHORT": true}
var validStatuses = map[string]bool{"OPEN": true, "CLOSED": true, "PENDING": true}

func CreateTrade(c *gin.Context) {
	userID := c.GetString("userID")
	accountID := c.GetString("accountID")

	var trade domains.Trade
	if err := c.ShouldBindJSON(&trade); err != nil {
		log.Printf("CreateTrade bind error: %v", err)
		response.BadRequest(c, "Datos inválidos")
		return
	}

	trade.UserID = userID
	trade.AccountID = accountID

	trade.Direction = strings.ToUpper(trade.Direction)
	trade.Status = strings.ToUpper(trade.Status)
	if trade.Status == "" {
		trade.Status = "CLOSED"
	}

	if trade.Symbol == "" {
		response.BadRequest(c, "El campo symbol es requerido")
		return
	}
	if len(trade.Symbol) > 20 {
		response.BadRequest(c, "symbol no puede superar 20 caracteres")
		return
	}
	if len(trade.Notes) > 5000 {
		response.BadRequest(c, "notes no puede superar 5000 caracteres")
		return
	}
	if !validDirections[trade.Direction] {
		response.BadRequest(c, "direction debe ser LONG o SHORT")
		return
	}
	if !validStatuses[trade.Status] {
		response.BadRequest(c, "status debe ser OPEN, CLOSED o PENDING")
		return
	}
	if trade.EntryPrice.IsZero() || trade.EntryPrice.IsNegative() {
		response.BadRequest(c, "entry_price debe ser mayor que 0")
		return
	}
	if trade.Size.IsZero() || trade.Size.IsNegative() {
		response.BadRequest(c, "size debe ser mayor que 0")
		return
	}
	if err := validateScreenshotPaths([]string(trade.Screenshots), userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	db := database.GetDB()
	if err := db.Create(&trade).Error; err != nil {
		log.Printf("CreateTrade db error accountID=%s: %v", accountID, err)
		response.InternalError(c)
		return
	}

	cache.Stats.Invalidate(accountID)
	c.JSON(http.StatusCreated, gin.H{"trade": trade})
}
