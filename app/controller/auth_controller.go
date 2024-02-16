package controller

import (
	"encoding/hex"
	"errors"
	"os"
	"time"

	"github.com/budgetin-app/user-service/app/constant"
	"github.com/budgetin-app/user-service/app/domain/model"
	"github.com/budgetin-app/user-service/app/pkg/hasher"
	"github.com/budgetin-app/user-service/app/pkg/mailer"
	"github.com/budgetin-app/user-service/app/repository"
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
	accountRepository           repository.AccountRepository
	loginInfoRepository         repository.LoginInfoRepository
	roleRepository              repository.RoleRepository
	sessionRepository           repository.SessionRepository
	emailVerificationRepository repository.EmailVerificationRepository
}

func NewAuthController(
	accountRepository repository.AccountRepository,
	loginInfoRepository repository.LoginInfoRepository,
	roleRepository repository.RoleRepository,
	sessionRepository repository.SessionRepository,
	emailVerificationRepository repository.EmailVerificationRepository,
) *AuthControllerImpl {
	return &AuthControllerImpl{
		accountRepository:           accountRepository,
		loginInfoRepository:         loginInfoRepository,
		roleRepository:              roleRepository,
		sessionRepository:           sessionRepository,
		emailVerificationRepository: emailVerificationRepository,
	}
}

func (c AuthControllerImpl) Register(username string, email string, password string) (*model.LoginInfo, error) {
	// Begin a transaction
	tx := c.accountRepository.BeginTransaction()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

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
	account := model.Account{RoleID: roleID}
	if err := tx.Create(&account).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Store the user credentials with the created account
	credential := model.LoginInfo{
		ID:            account.ID,
		Username:      username,
		Email:         email,
		PasswordHash:  string(hashedPassword),
		PasswordSalt:  hex.EncodeToString(passwordSalt),
		HashAlgorithm: model.HashAlgorithm{Name: string(hashAlgorithm)},
		EmailVerification: model.EmailVerification{
			Token:  uuid.New().String(),
			Status: model.VerificationPending,
		},
	}
	if err := tx.Create(&credential).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit the transaction if everything is successful
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Send email verification email asyncronously
	go c.sendVerificationEmail(&credential)

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
	// Find the email verification status
	credential, err := c.loginInfoRepository.FindLoginInfo(&model.LoginInfo{Email: email})
	if err != nil {
		log.WithError(err).Error("failed to find credential")
		return false, err
	}

	// Return email verification status
	verified := credential.EmailVerification.Status == model.EmailVerified

	// Send an email verification request to the target user when not yet verified, only sent
	// the email with a certain interval (minutes).
	resendInterval := 15 // TODO: Move the interval into service configuration
	if !verified && credential.EmailVerification.UpdatedAt.Add(time.Duration(resendInterval)*time.Minute).Before(time.Now()) {
		log.Debug("Send email")
		go c.sendVerificationEmail(&credential)
	} else {
		log.Debugf("Email already sent. Wait for %d minutes to resend", resendInterval)
	}

	return verified, nil
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

func (c *AuthControllerImpl) sendVerificationEmail(credential *model.LoginInfo) {
	if err := mailer.SendEmailVerification(
		credential.Email,
		credential.Username,
		credential.EmailVerification.Token,
		credential.EmailVerification.ExpiredAt,
	); err != nil {
		// Log error
		log.Errorf("error sending verification email: %v", err)
	} else {
		// Update the email verification status to 'sent'
		credential.EmailVerification.ID = credential.EmailVerificationID
		credential.EmailVerification.Status = model.VerificationSent
		c.emailVerificationRepository.UpdateEmailVerification(&credential.EmailVerification)
	}
}
