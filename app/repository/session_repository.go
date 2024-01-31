package repository

import (
	"log"

	"github.com/Budgetin-Project/user-service/app/domain/model"
	"gorm.io/gorm"
)

type SessionRepository interface {
	CreateSession(session *model.Session) (model.Session, error)
	FindSessionByToken(token string) (model.Session, error)
	UpdateSessionStatus(sessionID uint, status string) (bool, error)
	DeleteSession(session *model.Session) (bool, error)
}

type SessionRepositoryImpl struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepositoryImpl {
	return &SessionRepositoryImpl{db: db}
}

func (r SessionRepositoryImpl) CreateSession(session *model.Session) (model.Session, error) {
	if err := r.db.Create(&session).Error; err != nil {
		log.Fatalf("error create new session: %v", err)
		return model.Session{}, err
	}
	return *session, nil
}

func (r SessionRepositoryImpl) FindSessionByToken(token string) (model.Session, error) {
	session := model.Session{Token: token}
	if err := r.db.Find(&session).Error; err != nil {
		log.Fatalf("error find session: %v", err)
		return model.Session{}, err
	}
	return session, nil
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
