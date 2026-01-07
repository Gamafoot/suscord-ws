package storage

import (
	"suscord_ws/internal/domain/storage/database"

	"gorm.io/gorm"
)

type GormStorage struct {
	user    *userStorage
	chat    *chatStorage
	session *sessionStorage
}

func NewGormStorage(db *gorm.DB) *GormStorage {
	return &GormStorage{
		user:    NewUserStorage(db),
		chat:    NewChatStorage(db),
		session: NewSessionStorage(db),
	}
}

func (s *GormStorage) User() database.UserStorage {
	return s.user
}

func (s *GormStorage) Chat() database.ChatStorage {
	return s.chat
}

func (s *GormStorage) Session() database.SessionStorage {
	return s.session
}
