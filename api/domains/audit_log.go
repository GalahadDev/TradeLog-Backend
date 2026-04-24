package domains

import "time"

type AuditLog struct {
	ID         string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	AdminID    string    `gorm:"type:uuid;not null;index"                       json:"admin_id"`
	Action     string    `gorm:"not null"                                        json:"action"`      // ej. "PUT /api/admin/users/:id"
	TargetType string    `gorm:"not null"                                        json:"target_type"` // ej. "user"
	TargetID   string    `gorm:"type:uuid"                                       json:"target_id"`
	IP         string    `gorm:"not null"                                        json:"ip"`
	RequestID  string    `gorm:"not null"                                        json:"request_id"`
	CreatedAt  time.Time `                                                       json:"created_at"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
