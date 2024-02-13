package repository

import (
	"github.com/Budgetin-Project/user-management-service/config/database"
	"github.com/Budgetin-Project/user-service/app/domain/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type EmailVerificationRepository interface {
	UpdateEmailVerification(verification *model.EmailVerification) (model.EmailVerification, error)
	DeleteEmailVerification(verification *model.EmailVerification) (bool, error)
}

type EmailVerificationRepositoryImpl struct {
	db *gorm.DB
}

func NewEmailVerificationRepository(db *gorm.DB) *EmailVerificationRepositoryImpl {
	return &EmailVerificationRepositoryImpl{db: db}
}

func (r EmailVerificationRepositoryImpl) UpdateEmailVerification(verification *model.EmailVerification) (model.EmailVerification, error) {
	result := r.db.Model(&model.EmailVerification{ID: verification.ID}).Updates(&verification)
	if result.Error != nil {
		log.Fatalf("error update email verification: %v", result.Error)
		return model.EmailVerification{}, database.HandleErrorDB(result.Error)
	}
	return *verification, nil
}

func (r EmailVerificationRepositoryImpl) DeleteEmailVerification(verification *model.EmailVerification) (bool, error) {
	result := r.db.Delete(&verification)
	if result.Error != nil {
		log.Fatalf("error update email verification: %v", result.Error)
		return false, database.HandleErrorDB(result.Error)
	}
	return result.RowsAffected > 0, nil
}
