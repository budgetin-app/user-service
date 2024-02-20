package controller

import (
	"encoding/hex"
	"errors"
	"os"
	"time"

	"github.com/budgetin-app/user-service/app/constant"
	"github.com/budgetin-app/user-service/app/domain/model"
	"github.com/budgetin-app/user-service/app/pkg/hasher"
	"github.com/budgetin-app/user-service/app/pkg/helper/token"
	"github.com/budgetin-app/user-service/app/pkg/mailer"
	"github.com/budgetin-app/user-service/app/repository"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type AuthController interface {
	Register(username string, email string, password string) (*model.LoginInfo, error)
	Login(isEmail bool, identifier string, password string) (*model.Session, error)
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

func (c AuthControllerImpl) Login(isEmail bool, identifier string, password string) (*model.Session, error) {
	// Verify user's credential
	credential := &model.LoginInfo{}
	if isEmail {
		credential.Email = identifier
	} else {
		credential.Username = identifier
	}
	if err := c.loginInfoRepository.FindLoginInfo(credential); err != nil {
		log.WithError(err).Error("failed to find credential")
		return nil, err
	}

	// Validates user's password
	hash := hasher.New(hasher.HashAlgorithm(credential.HashAlgorithm.Name))
	salt, err := hex.DecodeString(credential.PasswordSalt)
	if err != nil {
		return nil, err
	}
	validPassword, err := hash.VerifyPassword(
		[]byte(credential.PasswordHash),
		[]byte(password),
		salt,
	)
	if err != nil {
		return nil, err
	} else if !validPassword {
		return nil, errors.New("password mismatched")
	}

	// Check for existing session
	oldSession, err := c.sessionRepository.FindActiveSession(credential.ID)
	if err != nil {
		log.Error(err)
	}
	if oldSession != nil {
		log.Debugf("found active session, id: %d", oldSession.ID)
		return nil, errors.New("user already logged on")
	}

	// Generate session token
	token, err := token.GenerateSessionToken()
	if err != nil {
		return nil, err
	}

	// Create new session for the user
	session, err := c.sessionRepository.CreateSession(credential.ID, token)
	if err != nil {
		return nil, err
	}

	// Return user session
	return &session, nil
}

func (c AuthControllerImpl) Logout(authToken string) (bool, error) {
	// Delete the session
	if err := c.sessionRepository.DeleteSessionByToken(authToken); err != nil {
		return false, err
	}
	return true, nil
}

func (c AuthControllerImpl) VerifyEmail(email string) (bool, error) {
	// Find the email verification status
	credential := &model.LoginInfo{Email: email}
	err := c.loginInfoRepository.FindLoginInfo(credential)
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
		go c.sendVerificationEmail(credential)
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
