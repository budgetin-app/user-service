package repository

import (
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/budgetin-app/user-service/app/domain/model"
	"gorm.io/gorm"
)

type RoleRepository interface {
	CreateRole(role *model.Role) (model.Role, error)
	AssignRolePermissions(role *model.Role, permissions ...model.Permission) error
	UpdateRole(newRole *model.Role) (model.Role, error)
	DeleteRole(role *model.Role) (bool, error)
}

type RoleRepositoryImpl struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepositoryImpl {
	return &RoleRepositoryImpl{db: db}
}

func (r RoleRepositoryImpl) CreateRole(role *model.Role) (model.Role, error) {
	if err := r.db.Create(&role).Error; err != nil {
		log.Errorf("error create new role: %v", err)
		return model.Role{}, err
	}
	return *role, nil
}

func (r RoleRepositoryImpl) AssignRolePermissions(role *model.Role, permissions ...model.Permission) error {
	// Check the permission ids
	if len(permissions) == 0 {
		return errors.New("permission id's should not be empty")
	}

	// Assign the permissions into the 'Permissions' field that representing the
	// many-to-many relationship
	role.Permissions = append(role.Permissions, permissions...)

	// Save the role with the updated permissions
	if err := r.db.Save(&role).Error; err != nil {
		log.Errorf("error save role: %v", err)
		return err
	}

	return nil
}

func (r RoleRepositoryImpl) UpdateRole(newRole *model.Role) (model.Role, error) {
	result := r.db.Model(&model.Role{ID: newRole.ID}).Updates(&newRole)
	if result.Error != nil {
		log.Errorf("error update role: %v", result.Error)
		return model.Role{}, result.Error
	}
	return *newRole, nil
}
func (r RoleRepositoryImpl) DeleteRole(role *model.Role) (bool, error) {
	result := r.db.Delete(&role)
	if result.Error != nil {
		log.Errorf("error delete role: %v", result.Error)
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}
