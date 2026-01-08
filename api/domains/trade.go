package domains

import (
	"time"

	"github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type Trade struct {
	ID     string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID string `gorm:"type:uuid;not null;index" json:"user_id"`

	Symbol    string `gorm:"not null" json:"symbol"`
	Direction string `gorm:"not null" json:"direction"`
	Status    string `gorm:"default:'CLOSED'" json:"status"`

	EntryPrice decimal.Decimal `gorm:"type:numeric;not null" json:"entry_price"`
	ExitPrice  decimal.Decimal `gorm:"type:numeric" json:"exit_price"`
	Size       decimal.Decimal `gorm:"type:numeric;not null" json:"size"`
	PnL        decimal.Decimal `gorm:"type:numeric;column:pnl" json:"pnl"`
	Commission decimal.Decimal `gorm:"type:numeric" json:"commission"`

	EntryDate time.Time  `gorm:"not null" json:"entry_date"`
	ExitDate  *time.Time `json:"exit_date"`

	Notes      string `json:"notes"`
	Screenshot string `gorm:"column:screenshot_url" json:"screenshot_url"`

	Tags pq.StringArray `gorm:"type:text[]" json:"tags"`

	CreatedAt time.Time `json:"created_at"`
}

func (Trade) TableName() string {
	return "trades"
}
