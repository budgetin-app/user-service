package controller

import (
	"errors"

	"github.com/Budgetin-Project/user-service/app/domain/model"
	"github.com/Budgetin-Project/user-service/app/repository"
)

type AuthController interface {
	Register(username string, email string, password string) (model.Account, error)
	Login(identifier string, isEmail bool, password string) (model.Session, error)
	Logout(authToken string) (bool, error)
	VerifyEmail(email string) (bool, error)
}

type AuthControllerImpl struct {
	accountRepository   repository.AccountRepository
	loginInfoRepository repository.LoginInfoRepository
	roleRepository      repository.RoleRepository
	sessionRepository   repository.SessionRepository
}

func NewAuthController(
	accountRepository repository.AccountRepository,
	loginInfoRepository repository.LoginInfoRepository,
	roleRepository repository.RoleRepository,
	sessionRepository repository.SessionRepository,
) *AuthControllerImpl {
	return &AuthControllerImpl{
		accountRepository:   accountRepository,
		loginInfoRepository: loginInfoRepository,
		roleRepository:      roleRepository,
		sessionRepository:   sessionRepository,
	}
}

func (c AuthControllerImpl) Register(username string, email string, password string) (model.Account, error) {
	// TODO: Not yet implemented
	return model.Account{}, errors.New("not yet implemented")
}

func (c AuthControllerImpl) Login(identifier string, isEmail bool, password string) (model.Session, error) {
	// TODO: Not yet implemented
	return model.Session{}, errors.New("not yet implemented")
}

func (c AuthControllerImpl) Logout(authToken string) (bool, error) {
	// TODO: Not yet implemented
	return false, errors.New("not yet implemented")
}

func (c AuthControllerImpl) VerifyEmail(email string) (bool, error) {
	// TODO: Not yet implemented
	return false, errors.New("not yet implemented")
}
