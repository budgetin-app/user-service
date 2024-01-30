package model

import "time"

type PasswordRecovery struct {
	ID        uint      `gorm:"column:password_recovery_id; primaryKey"`
	Token     string    `gorm:"column:recovery_token; size:100"`
	ExpiredAt time.Time `gorm:"column:token_expiration"`
	BaseModel
}

func (PasswordRecovery) TableName() string {
	return "password_recovery_info"
}
