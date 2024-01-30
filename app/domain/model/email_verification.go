package model

import "time"

type EmailVerification struct {
	ID        uint      `gorm:"column:email_verification_id; primaryKey"`
	Status    string    `gorm:"size:50"`
	Token     string    `gorm:"column:verification_token; size:100; unique"`
	ExpiredAt time.Time `gorm:"column:token_expiration"`
	BaseModel
}

func (EmailVerification) TableName() string {
	return "email_verification_info"
}
