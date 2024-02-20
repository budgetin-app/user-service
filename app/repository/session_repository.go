package repository

import (
	"errors"
	"time"

	"github.com/budgetin-app/user-service/app/domain/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SessionRepository interface {
	CreateSession(userID uint, token string) (model.Session, error)
	FindActiveSession(userID uint) (*model.Session, error)
	UpdateSessionStatus(sessionID uint, status string) (bool, error)
	DeleteSessionByToken(authToken string) error
}

type SessionRepositoryImpl struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepositoryImpl {
	return &SessionRepositoryImpl{db: db}
}

func (r SessionRepositoryImpl) CreateSession(userID uint, token string) (model.Session, error) {
	session := &model.Session{
		UserID: userID,
		Token:  token,
	}
	if err := r.db.Create(&session).Error; err != nil {
		log.Errorf("error create new session: %v", err)
		return model.Session{}, err
	}
	return *session, nil
}

func (r SessionRepositoryImpl) FindActiveSession(userID uint) (*model.Session, error) {
	var session model.Session

	// Find the last active session associated with the given userID
	if err := r.db.Where("user_id = ? AND session_expiration > ?", userID, time.Now()).
		Order("session_expiration desc").First(&session).Error; err != nil {
		return nil, err
	}

	return &session, nil
}

func (r SessionRepositoryImpl) UpdateSessionStatus(sessionID uint, status string) (bool, error) {
	result := r.db.Model(model.Session{ID: sessionID}).Update("status", status)
	if result.Error != nil {
		log.Errorf("error finish session: %v", result.Error)
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r SessionRepositoryImpl) DeleteSessionByToken(authToken string) error {
	result := r.db.Where("session_token = ?", authToken).Delete(&model.Session{})
	if result.Error != nil {
		log.Errorf("error delete session: %v", result.Error)
		return result.Error
	}

	log.Debugf("session deleted count: %d", result.RowsAffected)
	if result.RowsAffected <= 0 {
		return errors.New("no session deleted")
	}

	return nil
}
