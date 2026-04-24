package accounts

import (
	"log"
	"net/http"
	"strings"

	"samll-trading-back/api/database"
	"samll-trading-back/api/domains"
	"samll-trading-back/api/response"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

var validAccountTypes = map[string]bool{"prop_firm": true, "personal": true}
var validAccountStatuses = map[string]bool{"active": true, "completed": true, "lost": true}

type createAccountReq struct {
	Name             string           `json:"name"`
	Broker           string           `json:"broker"`
	AccountType      string           `json:"account_type"`
	InitialBalance   decimal.Decimal  `json:"initial_balance"`
	ProfitTarget     *decimal.Decimal `json:"profit_target"`
	MaxDrawdownLimit *decimal.Decimal `json:"max_drawdown_limit"`
	Currency         string           `json:"currency"`
}

func CreateAccount(c *gin.Context) {
	userID := c.GetString("userID")

	var req createAccountReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("CreateAccount bind error: %v", err)
		response.BadRequest(c, "Datos inválidos")
		return
	}

	req.AccountType = strings.ToLower(req.AccountType)

	if req.Name == "" {
		response.BadRequest(c, "El campo name es requerido")
		return
	}
	if req.Broker == "" {
		response.BadRequest(c, "El campo broker es requerido")
		return
	}
	if !validAccountTypes[req.AccountType] {
		response.BadRequest(c, "account_type debe ser prop_firm o personal")
		return
	}
	if req.InitialBalance.IsZero() || req.InitialBalance.IsNegative() {
		response.BadRequest(c, "initial_balance debe ser mayor que 0")
		return
	}
	if req.AccountType == "personal" && (req.ProfitTarget != nil || req.MaxDrawdownLimit != nil) {
		response.BadRequest(c, "profit_target y max_drawdown_limit solo aplican a cuentas prop_firm")
		return
	}
	if req.ProfitTarget != nil && (req.ProfitTarget.IsZero() || req.ProfitTarget.IsNegative()) {
		response.BadRequest(c, "profit_target debe ser mayor que 0")
		return
	}
	if req.MaxDrawdownLimit != nil && (req.MaxDrawdownLimit.IsZero() || req.MaxDrawdownLimit.IsNegative()) {
		response.BadRequest(c, "max_drawdown_limit debe ser mayor que 0")
		return
	}

	currency := req.Currency
	if currency == "" {
		currency = "USD"
	}

	account := domains.TradingAccount{
		UserID:           userID,
		Name:             req.Name,
		Broker:           req.Broker,
		AccountType:      req.AccountType,
		Status:           "active",
		Currency:         currency,
		InitialBalance:   req.InitialBalance,
		ProfitTarget:     req.ProfitTarget,
		MaxDrawdownLimit: req.MaxDrawdownLimit,
	}

	db := database.GetDB()
	if err := db.Create(&account).Error; err != nil {
		log.Printf("CreateAccount db error userID=%s: %v", userID, err)
		response.InternalError(c)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"account": account})
}
