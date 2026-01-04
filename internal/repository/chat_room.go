package repository

import (
	"github.com/MrBista/blog-api/internal/exception"
	"github.com/MrBista/blog-api/internal/models"
	"gorm.io/gorm"
)

type ChatRoomRepository interface {
	CreateRoom(*models.ChatRoom) error
	FindByBlogAndReader(blogId, readerId int) (*models.ChatRoom, error)
}

type ChatRoomImpl struct {
	DB *gorm.DB
}

func NewChatRoom(db *gorm.DB) ChatRoomRepository {
	return &ChatRoomImpl{
		DB: db,
	}
}

func (r *ChatRoomImpl) CreateRoom(data *models.ChatRoom) error {
	if err := r.DB.Create(data).Error; err != nil {
		return exception.NewGormDBErr(err)
	}
	return nil
}

func (r *ChatRoomImpl) FindByBlogAndReader(blogId, readerId int) (*models.ChatRoom, error) {

	var room *models.ChatRoom

	err := r.DB.Where("post_id = ? And reader_id = ?", blogId, readerId).First(&room).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exception.NewNotFoundErr("room not found")
		}
		return nil, exception.NewGormDBErr(err)
	}

	return room, nil
}
