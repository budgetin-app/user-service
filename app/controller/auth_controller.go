package controller

import (
	"encoding/hex"
	"errors"
	"os"

	"github.com/Budgetin-Project/user-service/app/constant"
	"github.com/Budgetin-Project/user-service/app/domain/model"
	"github.com/Budgetin-Project/user-service/app/pkg/hasher"
	"github.com/Budgetin-Project/user-service/app/pkg/mailer"
	"github.com/Budgetin-Project/user-service/app/repository"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type AuthController interface {
	Register(username string, email string, password string) (*model.LoginInfo, error)
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

func (c AuthControllerImpl) Register(username string, email string, password string) (*model.LoginInfo, error) {
	// Generate hashed password with random salt
	hashAlgorithm := getHashAlgorithm()
	hash := hasher.New(hashAlgorithm)
	passwordSalt := hasher.GenerateRandomSalt()
	hashedPassword, err := hash.GenerateHashPassword([]byte(password), passwordSalt)
	if err != nil {
		return nil, err
	}

	// TODO: Later the roleID should be retrieved from the request parameter
	// Required because the 'Register' method can be used to register user dynamically.
	//
	// There will be, role level check (may need to add new field or just use the id
	// to define the role level)
	//
	// For example, to create user 'Admin'
	roleID := constant.UserRoleID

	// Store the user account
	account, err := c.accountRepository.CreateAccount(&model.Account{RoleID: roleID})
	if err != nil {
		return nil, err
	}

	// Store the user credentials with the created account
	credential, err := c.loginInfoRepository.CreateLoginInfo(&model.LoginInfo{
		ID:           account.ID,
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		PasswordSalt: hex.EncodeToString(passwordSalt),
		HashAlgorithm: model.HashAlgorithm{
			Name: string(hashAlgorithm),
		},
		EmailVerification: model.EmailVerification{
			Token:  uuid.New().String(),
			Status: model.VerificationPending,
		},
	})
	if err != nil {
		// Delete the account if error
		c.accountRepository.DeleteAccount(&account)
		return nil, err
	}

	// Send verification email using event bus to trigger sending email verification asyncronously
	go func() {
		if err := mailer.SendEmailVerification(
			credential.Email,
			credential.Username,
			credential.EmailVerification.Token,
			credential.EmailVerification.ExpiredAt,
		); err != nil {
			// Log error
			log.Errorf("error sending verification email: %v", err)
		} else {
			// TODO: Update the email verification status to 'sent'
		}
	}()

	return &credential, nil
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

func getHashAlgorithm() hasher.HashAlgorithm {
	// Use 'bcrypt' as the default hashing algorithm
	algorithm := hasher.BCrypt

	val := os.Getenv("PASSWORD_HASH_ALGORITHM")
	if val != "" && hasher.IsAlgorithmAllowed(hasher.HashAlgorithm(val)) {
		algorithm = hasher.HashAlgorithm(val)
	}

	return algorithm
}
