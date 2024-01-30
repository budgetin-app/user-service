package model

type Session struct {
	ID        uint `gorm:"column:session_id; primaryKey"`
	UserID    uint
	User      Account
	Token     string `gorm:"column:session_token; size:100; unique"`
	ExpiredAt string `gorm:"column:session_expiration; size:100"`
	BaseModel
}

func (Session) TableName() string {
	return "user_sessions"
}
