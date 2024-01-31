package repository

import (
	"log"

	"github.com/Budgetin-Project/user-service/app/domain/model"
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
	if err := r.db.Create(&info).Error; err != nil {
		log.Fatalf("error create new login info: %v", err)
		return model.LoginInfo{}, err
	}
	return *info, nil
}

func (r LoginInfoRepositoryImpl) FindLoginInfo(info *model.LoginInfo) (model.LoginInfo, error) {
	if err := r.db.Find(&info).Error; err != nil {
		log.Fatalf("error find login info: %v", err)
		return model.LoginInfo{}, err
	}
	return *info, nil
}

func (r LoginInfoRepositoryImpl) UpdateLoginInfo(newInfo *model.LoginInfo) (model.LoginInfo, error) {
	result := r.db.Model(&model.LoginInfo{ID: newInfo.ID}).Updates(&newInfo)
	if result.Error != nil {
		log.Fatalf("error update login info: %v", result.Error)
		return model.LoginInfo{}, result.Error
	}
	return *newInfo, nil
}

func (r LoginInfoRepositoryImpl) DeleteLoginInfo(info *model.LoginInfo) (bool, error) {
	result := r.db.Delete(&info)
	if result.Error != nil {
		log.Fatalf("error delete login info: %v", result.Error)
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}
