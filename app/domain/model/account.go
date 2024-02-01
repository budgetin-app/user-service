package model

import "time"

type Account struct {
	ID          uint    `gorm:"column:user_id; primaryKey"`
	UserName    *string `gorm:"size:100"`
	Gender      *string `gorm:"size:1"`
	DateOfBirth time.Time
	RoleID      uint
	Role        Role
	BaseModel
}

func (Account) TableName() string {
	return "user_accounts"
}
