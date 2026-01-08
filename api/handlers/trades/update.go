package trades

import (
	"net/http"
	"strings"
	"time"

	"samll-trading-back/api/database"
	"samll-trading-back/api/domains"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type TradeUpdateReq struct {
	// Datos editables
	Symbol     *string          `json:"symbol"`
	Direction  *string          `json:"direction"`
	Status     *string          `json:"status"`
	EntryPrice *decimal.Decimal `json:"entry_price"`
	ExitPrice  *decimal.Decimal `json:"exit_price"`
	Size       *decimal.Decimal `json:"size"`
	PnL        *decimal.Decimal `json:"pnl"`
	Commission *decimal.Decimal `json:"commission"`

	// Fechas y Extras
	EntryDate  *time.Time `json:"entry_date"`
	ExitDate   *time.Time `json:"exit_date"`
	Notes      *string    `json:"notes"`
	Screenshot *string    `json:"screenshot_url"`
	Tags       []string   `json:"tags"`
}

func UpdateTrade(c *gin.Context) {
	userID := c.GetString("userID")
	tradeID := c.Param("id")

	var req TradeUpdateReq
	// Validar JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
		return
	}

	db := database.GetDB()
	var trade domains.Trade

	// 1. Verificar propiedad del trade antes de editar
	if err := db.Where("id = ? AND user_id = ?", tradeID, userID).First(&trade).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Trade no encontrado"})
		return
	}

	// 2. Construir mapa de actualización dinámica
	updates := make(map[string]interface{})

	if req.Symbol != nil {
		updates["symbol"] = *req.Symbol
	}
	if req.Direction != nil {
		upper := strings.ToUpper(*req.Direction)
		updates["direction"] = upper
	}
	if req.Status != nil {
		upper := strings.ToUpper(*req.Status)
		updates["status"] = upper
	}

	// Decimales
	if req.EntryPrice != nil {
		updates["entry_price"] = *req.EntryPrice
	}
	if req.ExitPrice != nil {
		updates["exit_price"] = *req.ExitPrice
	}
	if req.Size != nil {
		updates["size"] = *req.Size
	}
	if req.PnL != nil {
		updates["pnl"] = *req.PnL
	}
	if req.Commission != nil {
		updates["commission"] = *req.Commission
	}

	// Fechas y Textos
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
		updates["screenshot_url"] = *req.Screenshot
	}

	// Arrays (Tags)
	if req.Tags != nil {
		updates["tags"] = req.Tags
	}

	// 3. Ejecutar Update
	if err := db.Model(&trade).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar trade"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Operación actualizada",
		"trade":   trade,
	})
}
