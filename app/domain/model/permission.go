package model

type Permission struct {
	ID   uint   `gorm:"column:permission_id; primaryKey"`
	Name string `gorm:"column:permission_name; size:50"`
	BaseModel
}
