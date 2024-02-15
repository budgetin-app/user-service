package repository

import (
	"github.com/budgetin-app/user-management-service/config/database"
	"github.com/budgetin-app/user-service/app/domain/model"
	"gorm.io/gorm"
)

type LoginInfoRepository interface {
	CreateLoginInfo(info *model.LoginInfo) (model.LoginInfo, error)
	FindLoginInfo(info *model.LoginInfo) (model.LoginInfo, error)
	UpdateLoginInfo(newInfo *model.LoginInfo) (model.LoginInfo, error)
	DeleteLoginInfo(info *model.LoginInfo) (bool, error)
}

type LoginInfoRepositoryImpl struct {
	db *gorm.DB
}

func NewLoginInfoRepository(db *gorm.DB) *LoginInfoRepositoryImpl {
	return &LoginInfoRepositoryImpl{db: db}
}

func (r LoginInfoRepositoryImpl) CreateLoginInfo(info *model.LoginInfo) (model.LoginInfo, error) {
	// Check hash algorithm already exists
	var hashAlgorithm model.HashAlgorithm
	if err := r.db.FirstOrCreate(&hashAlgorithm, &info.HashAlgorithm).Error; err != nil {
		return model.LoginInfo{}, database.HandleErrorDB(err)
	}

	// Create login info using the hashAlgorithm found
	info.HashAlgorithm = hashAlgorithm
	if err := r.db.Create(&info).Error; err != nil {
		return model.LoginInfo{}, database.HandleErrorDB(err)
	}
	return *info, nil
}

func (r LoginInfoRepositoryImpl) FindLoginInfo(info *model.LoginInfo) (model.LoginInfo, error) {
	if err := r.db.Find(&info).Error; err != nil {
		return model.LoginInfo{}, database.HandleErrorDB(err)
	}
	return *info, nil
}

func (r LoginInfoRepositoryImpl) UpdateLoginInfo(newInfo *model.LoginInfo) (model.LoginInfo, error) {
	result := r.db.Model(&model.LoginInfo{ID: newInfo.ID}).Updates(&newInfo)
	if result.Error != nil {
		return model.LoginInfo{}, database.HandleErrorDB(result.Error)
	}
	return *newInfo, nil
}

func (r LoginInfoRepositoryImpl) DeleteLoginInfo(info *model.LoginInfo) (bool, error) {
	result := r.db.Delete(&info)
	if result.Error != nil {
		return false, database.HandleErrorDB(result.Error)
	}
	return result.RowsAffected > 0, nil
}
