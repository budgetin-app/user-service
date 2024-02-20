package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	SessionExpDuration = 1 // hours
)

type Session struct {
	ID        uint `gorm:"column:session_id; primaryKey"`
	UserID    uint
	User      Account
	Token     string    `gorm:"column:session_token; size:100; unique"`
	ExpiredAt time.Time `gorm:"column:session_expiration"`
	BaseModel
}

func (Session) TableName() string {
	return "user_sessions"
}

func (s *Session) BeforeCreate(tx *gorm.DB) (err error) {
	s.ExpiredAt = time.Now().Add(time.Hour * SessionExpDuration)
	return
}
