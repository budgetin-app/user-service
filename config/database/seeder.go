package database

import (
	"github.com/budgetin-app/user-service/app/constant"
	"github.com/budgetin-app/user-service/app/domain/model"
	"gorm.io/gorm"
)

// SeederDB seeds the database with the defined values
// Note: Use 'Save' method when seeding the db, so it doesn't create duplicates
func SeederDB(db *gorm.DB) {
	// Seed roles
	for roleID, roleName := range constant.GetUserRoles() {
		if err := db.Save(&model.Role{ID: roleID, Name: roleName}).Error; err != nil {
			panic("failed to seed roles")
		}
	}

	// .. add db seeder here
}
