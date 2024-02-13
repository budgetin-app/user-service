package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	VerificationPending = "pending"
	VerificationSent    = "sent"
	VerificationError   = "error"
	EmailVerified       = "verified"

	// TODO: Should moved to configuration database
	TokenExpDuration = 24 // hours
)

type EmailVerification struct {
	ID        uint      `gorm:"column:email_verification_id; primaryKey"`
	Status    string    `gorm:"size:50; default:pending"`
	Token     string    `gorm:"column:verification_token; size:100; unique"`
	ExpiredAt time.Time `gorm:"column:token_expiration"`
	BaseModel
}

func (EmailVerification) TableName() string {
	return "email_verification_info"
}

func (e *EmailVerification) BeforeCreate(tx *gorm.DB) (err error) {
	e.setTokenExpirationTime()
	return
}

func (e *EmailVerification) BeforeSave(tx *gorm.DB) (err error) {
	e.setTokenExpirationTime()
	return
}

// setTokenExpirationTime sets the token expiration time for
func (e *EmailVerification) setTokenExpirationTime() {
	switch e.Status {
	case VerificationPending, VerificationError:
		e.ExpiredAt = time.Now().Add(time.Hour * TokenExpDuration)
	default:
		// Do nothing when status is 'sent' or 'verified'
	}
}
