package repository

import (
	"errors"
	"fmt"

	"github.com/MrBista/blog-api/internal/models"
	"gorm.io/gorm"
)

type PostRepository interface {
	GetAllPost() ([]models.Post, error)
	GetDetailPost(slug string) (*models.Post, error)
	CreatePost(post *models.Post) error
	UpdatePost(id int, data map[string]interface{}) error
	DeletePost(id int) error
}

type PostRepositoryImpl struct {
	DB *gorm.DB
}

func NewPostRepository(DB *gorm.DB) PostRepository {
	return &PostRepositoryImpl{
		DB: DB,
	}
}

func (r *PostRepositoryImpl) GetAllPost() ([]models.Post, error) {
	var posts []models.Post
	tx := r.DB.Find(&posts)

	return posts, tx.Error
}

func (r *PostRepositoryImpl) GetDetailPost(slug string) (*models.Post, error) {
	var post models.Post
	tx := r.DB.Take(&post, "slug like ", "%"+slug+"%")

	if tx.Error != nil {
		return nil, fmt.Errorf("failed to get post %w", tx.Error)
	}

	return &post, tx.Error

}

func (r *PostRepositoryImpl) CreatePost(post *models.Post) error {
	txRes := r.DB.Create(post)

	if txRes.Error != nil {
		return fmt.Errorf("failed to save post %w", txRes.Error)
	}

	return txRes.Error
}

func (r *PostRepositoryImpl) UpdatePost(id int, data map[string]interface{}) error {
	res := r.DB.Where("id = ?", id).Updates(data)

	if res.RowsAffected == 0 {
		return errors.New("no row affacted")
	}

	return fmt.Errorf("")
}

func (r *PostRepositoryImpl) DeletePost(id int) error {
	rxRes := r.DB.Where("id = ?", id).Delete(&models.User{})

	return rxRes.Error
}
