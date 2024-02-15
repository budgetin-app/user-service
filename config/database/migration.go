package database

import (
	"github.com/budgetin-app/user-service/app/domain/model"
	"gorm.io/gorm"
)

// MigrateDB execute the database migration according to the models
func MigrateDB(db *gorm.DB) {
	db.AutoMigrate(
		&model.Account{},
		&model.Role{},
		&model.Permission{},
		&model.Session{},
		&model.LoginInfo{},
		&model.HashAlgorithm{},
		&model.EmailVerification{},
		&model.PasswordRecovery{},
		// .. add other db migration model here
	)
}
