package repository

import (
	"log"

	"github.com/budgetin-app/user-service/app/domain/model"
	"gorm.io/gorm"
)

type AccountRepository interface {
	CreateAccount(account *model.Account) (model.Account, error)
	FindAccountByUserID(userID uint) (model.Account, error)
	UpdateAccount(newAccount *model.Account) (model.Account, error)
	DeleteAccount(account *model.Account) (bool, error)
	BeginTransaction() *gorm.DB
}

type AccountRepositoryImpl struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepositoryImpl {
	return &AccountRepositoryImpl{db: db}
}

func (r AccountRepositoryImpl) CreateAccount(account *model.Account) (model.Account, error) {
	if err := r.db.Create(&account).Error; err != nil {
		log.Fatalf("error create new account: %v", err)
		return model.Account{}, err
	}
	return *account, nil
}

func (r AccountRepositoryImpl) FindAccountByUserID(userID uint) (model.Account, error) {
	account := model.Account{ID: userID}
	if err := r.db.Find(&account).Error; err != nil {
		log.Fatalf("error find account by user id: %v", err)
		return model.Account{}, err
	}
	return account, nil
}

func (r AccountRepositoryImpl) UpdateAccount(newAccount *model.Account) (model.Account, error) {
	result := r.db.Model(&model.Account{ID: newAccount.ID}).Updates(&newAccount)
	if result.Error != nil {
		log.Fatalf("error update account: %v", result.Error)
		return model.Account{}, result.Error
	}
	return *newAccount, nil
}

func (r AccountRepositoryImpl) DeleteAccount(account *model.Account) (bool, error) {
	result := r.db.Delete(&account)
	if result.Error != nil {
		log.Fatalf("error delete account: %v", result.Error)
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r AccountRepositoryImpl) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}
