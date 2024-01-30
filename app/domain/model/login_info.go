package model

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
