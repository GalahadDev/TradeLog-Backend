package trades

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"samll-trading-back/api/cache"
	"samll-trading-back/api/database"
	"samll-trading-back/api/domains"
	"samll-trading-back/api/response"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type TradeUpdateReq struct {
	Symbol     *string          `json:"symbol"`
	Direction  *string          `json:"direction"`
	Status     *string          `json:"status"`
	EntryPrice *decimal.Decimal `json:"entry_price"`
	ExitPrice  *decimal.Decimal `json:"exit_price"`
	Size       *decimal.Decimal `json:"size"`
	PnL        *decimal.Decimal `json:"pnl"`
	Commission *decimal.Decimal `json:"commission"`

	EntryDate  *time.Time `json:"entry_date"`
	ExitDate   *time.Time `json:"exit_date"`
	Notes      *string    `json:"notes"`
	Screenshot *string    `json:"screenshot_url"`
	Tags       []string   `json:"tags"`
}

func UpdateTrade(c *gin.Context) {
	accountID := c.GetString("accountID")
	tradeID := c.Param("trade_id")

	var req TradeUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("UpdateTrade bind error: %v", err)
		response.BadRequest(c, "Datos inválidos")
		return
	}

	db := database.GetDB()
	var trade domains.Trade

	if err := db.Where("id = ? AND account_id = ?", tradeID, accountID).First(&trade).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "Trade no encontrado")
		} else {
			log.Printf("UpdateTrade fetch error tradeID=%s accountID=%s: %v", tradeID, accountID, err)
			response.InternalError(c)
		}
		return
	}

	updates := make(map[string]any)

	if req.Symbol != nil {
		updates["symbol"] = *req.Symbol
	}
	if req.Direction != nil {
		upper := strings.ToUpper(*req.Direction)
		if !validDirections[upper] {
			response.BadRequest(c, "direction debe ser LONG o SHORT")
			return
		}
		updates["direction"] = upper
	}
	if req.Status != nil {
		upper := strings.ToUpper(*req.Status)
		if !validStatuses[upper] {
			response.BadRequest(c, "status debe ser OPEN, CLOSED o PENDING")
			return
		}
		updates["status"] = upper
	}
	if req.EntryPrice != nil {
		if req.EntryPrice.IsZero() || req.EntryPrice.IsNegative() {
			response.BadRequest(c, "entry_price debe ser mayor que 0")
			return
		}
		updates["entry_price"] = *req.EntryPrice
	}
	if req.ExitPrice != nil {
		updates["exit_price"] = *req.ExitPrice
	}
	if req.Size != nil {
		if req.Size.IsZero() || req.Size.IsNegative() {
			response.BadRequest(c, "size debe ser mayor que 0")
			return
		}
		updates["size"] = *req.Size
	}
	if req.PnL != nil {
		updates["pnl"] = *req.PnL
	}
	if req.Commission != nil {
		updates["commission"] = *req.Commission
	}
	if req.EntryDate != nil {
		updates["entry_date"] = *req.EntryDate
	}
	if req.ExitDate != nil {
		updates["exit_date"] = *req.ExitDate
	}
	if req.Notes != nil {
		updates["notes"] = *req.Notes
	}
	if req.Screenshot != nil {
		if *req.Screenshot != "" {
			parsed, err := url.Parse(*req.Screenshot)
			if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
				response.BadRequest(c, "screenshot_url debe ser una URL válida con esquema http o https")
				return
			}
		}
		updates["screenshot_url"] = *req.Screenshot
	}
	if req.Tags != nil {
		updates["tags"] = pq.StringArray(req.Tags)
	}

	if err := db.Model(&trade).Updates(updates).Error; err != nil {
		log.Printf("UpdateTrade db error tradeID=%s accountID=%s: %v", tradeID, accountID, err)
		response.InternalError(c)
		return
	}

	cache.Stats.Invalidate(accountID)
	c.JSON(http.StatusOK, gin.H{"message": "Operación actualizada", "trade": trade})
}
