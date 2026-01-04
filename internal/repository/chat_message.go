package repository

import (
	"github.com/MrBista/blog-api/internal/exception"
	"github.com/MrBista/blog-api/internal/models"
	"gorm.io/gorm"
)

type ChatMessageRepository interface {
	SaveMessage(chatDetail *models.ChatMessage) error
}

type ChatMessageRepositoryImpl struct {
	DB *gorm.DB
}

func NewChatMessageRepository(db *gorm.DB) ChatMessageRepository {

	return &ChatMessageRepositoryImpl{
		DB: db,
	}
}

func (r *ChatMessageRepositoryImpl) SaveMessage(chatDetail *models.ChatMessage) error {

	if err := r.DB.Create(chatDetail).Error; err != nil {
		return exception.NewGormDBErr(err)
	}

	return nil
}
