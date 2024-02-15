package model

import (
	"errors"

	"gorm.io/gorm"
)

type LoginInfo struct {
	ID                  uint   `gorm:"column:user_id; primaryKey"`
	Username            string `gorm:"size:20; unique"`
	Email               string `gorm:"size:100; unique"`
	PasswordHash        string `gorm:"size:250"`
	PasswordSalt        string `gorm:"size:100"`
	HashAlgorithmID     uint
	HashAlgorithm       HashAlgorithm `gorm:"foreignKey:HashAlgorithmID; references:ID"`
	EmailVerificationID uint
	EmailVerification   EmailVerification `gorm:"foreignKey:EmailVerificationID; references:ID"`
	PasswordRecoveryID  *uint
	PasswordRecovery    PasswordRecovery `gorm:"foreignKey:PasswordRecoveryID; references:ID"`
	BaseModel
}

func (LoginInfo) TableName() string {
	return "user_login_info"
}

func (i *LoginInfo) BeforeCreate(tx *gorm.DB) (err error) {
	// Checking username or email is already exists
	var info LoginInfo
	if res := tx.Where("username = ?", i.Username).Or("email = ?", i.Email).Find(&info); res.RowsAffected > 0 {
		if i.Username == info.Username {
			return errors.New("username already exists")
		}
		if i.Email == info.Email {
			return errors.New("email already exists")
		}

	}
	return nil
}

func (i *LoginInfo) BeforeUpdate(tx *gorm.DB) (err error) {
	// Username should not be changed at any circumstances, because it will be
	// used in the password hashing process
	if tx.Statement.Changed("Username") {
		return errors.New("username not allowed to be changed")
	}
	return nil
}
