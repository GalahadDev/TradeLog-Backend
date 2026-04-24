package domains

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type TradingAccount struct {
	ID     string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID string `gorm:"type:uuid;not null;index"                       json:"user_id"`

	Name        string `gorm:"not null"        json:"name"`
	Broker      string `gorm:"not null"        json:"broker"`
	AccountType string `gorm:"not null"        json:"account_type"`
	Status      string `gorm:"default:'active'" json:"status"`
	Currency    string `gorm:"default:'USD'"   json:"currency"`

	InitialBalance   decimal.Decimal  `gorm:"type:numeric;not null" json:"initial_balance"`
	ProfitTarget     *decimal.Decimal `gorm:"type:numeric"          json:"profit_target,omitempty"`
	MaxDrawdownLimit *decimal.Decimal `gorm:"type:numeric"          json:"max_drawdown_limit,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (TradingAccount) TableName() string {
	return "trading_accounts"
}
