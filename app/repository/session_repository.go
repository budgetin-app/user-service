package repository

import (
	"log"
	"time"

	"github.com/budgetin-app/user-service/app/domain/model"
	"gorm.io/gorm"
)

type SessionRepository interface {
	CreateSession(userID uint, token string) (model.Session, error)
	FindActiveSession(userID uint) (*model.Session, error)
	UpdateSessionStatus(sessionID uint, status string) (bool, error)
	DeleteSession(session *model.Session) (bool, error)
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
		log.Fatalf("error create new session: %v", err)
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
		log.Fatalf("error finish session: %v", result.Error)
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r SessionRepositoryImpl) DeleteSession(session *model.Session) (bool, error) {
	result := r.db.Delete(&session)
	if result.Error != nil {
		log.Fatalf("error delete session: %v", result.Error)
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}
