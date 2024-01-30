package model

type LoginExternal struct {
	ID         uint             `gorm:"column:user_id; primaryKey"`
	ProviderID uint             `gorm:"column:external_provider_id"`
	Provider   ExternalProvider `gorm:"foreignKey:ProviderID; references:ID"`
	Token      string           `gorm:"column:provider_token; size:100"`
	BaseModel
}
