package model

type ExternalProvider struct {
	ID   uint   `gorm:"column:external_provider_id; primaryKey"`
	Name string `gorm:"column:provider_name; size:50; unique"`
	Url  string `gorm:"column:web_service_url"`
	BaseModel
}
