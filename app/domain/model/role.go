package model

type Role struct {
	ID          uint         `gorm:"column:role_id; primaryKey"`
	Name        string       `gorm:"column:role_name; size:20"`
	Permissions []Permission `gorm:"many2many:granted_permissions; foreignKey:ID; joinForeignKey:RoleID; Reference:PermissionID; joinReference:PermissionID"`
	BaseModel
}

func (Role) TableName() string {
	return "user_roles"
}
