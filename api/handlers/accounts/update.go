package accounts

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"samll-trading-back/api/database"
	"samll-trading-back/api/domains"
	"samll-trading-back/api/response"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type accountUpdateReq struct {
	Name             *string          `json:"name"`
	Broker           *string          `json:"broker"`
	Status           *string          `json:"status"`
	ProfitTarget     *decimal.Decimal `json:"profit_target"`
	MaxDrawdownLimit *decimal.Decimal `json:"max_drawdown_limit"`
}

func UpdateAccount(c *gin.Context) {
	accountID := c.GetString("accountID")

	var req accountUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("UpdateAccount bind error: %v", err)
		response.BadRequest(c, "Datos inválidos")
		return
	}

	db := database.GetDB()
	var account domains.TradingAccount

	if err := db.First(&account, "id = ?", accountID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "Cuenta no encontrada")
		} else {
			log.Printf("UpdateAccount fetch error accountID=%s: %v", accountID, err)
			response.InternalError(c)
		}
		return
	}

	updates := make(map[string]any)

	if req.Name != nil {
		if strings.TrimSpace(*req.Name) == "" {
			response.BadRequest(c, "name no puede estar vacío")
			return
		}
		updates["name"] = *req.Name
	}
	if req.Broker != nil {
		if strings.TrimSpace(*req.Broker) == "" {
			response.BadRequest(c, "broker no puede estar vacío")
			return
		}
		updates["broker"] = *req.Broker
	}
	if req.Status != nil {
		lower := strings.ToLower(*req.Status)
		if !validAccountStatuses[lower] {
			response.BadRequest(c, "status debe ser active, completed o lost")
			return
		}
		updates["status"] = lower
	}
	if req.ProfitTarget != nil {
		if account.AccountType != "prop_firm" {
			response.BadRequest(c, "profit_target solo aplica a cuentas prop_firm")
			return
		}
		if req.ProfitTarget.IsZero() || req.ProfitTarget.IsNegative() {
			response.BadRequest(c, "profit_target debe ser mayor que 0")
			return
		}
		updates["profit_target"] = *req.ProfitTarget
	}
	if req.MaxDrawdownLimit != nil {
		if account.AccountType != "prop_firm" {
			response.BadRequest(c, "max_drawdown_limit solo aplica a cuentas prop_firm")
			return
		}
		if req.MaxDrawdownLimit.IsZero() || req.MaxDrawdownLimit.IsNegative() {
			response.BadRequest(c, "max_drawdown_limit debe ser mayor que 0")
			return
		}
		updates["max_drawdown_limit"] = *req.MaxDrawdownLimit
	}

	if len(updates) == 0 {
		c.JSON(http.StatusOK, gin.H{"account": account})
		return
	}

	if err := db.Model(&account).Updates(updates).Error; err != nil {
		log.Printf("UpdateAccount db error accountID=%s: %v", accountID, err)
		response.InternalError(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"account": account})
}
