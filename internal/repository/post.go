package repository

import (
	"github.com/MrBista/blog-api/internal/models"
	"gorm.io/gorm"
)

type PostRepository interface {
	GetAllPost() ([]models.Post, error)
	GetDetailPost(slug string) (models.Post, error)
	CreatePost() error
	UpdatePost() error
	DeletePost() error
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

func (r *PostRepositoryImpl) GetDetailPost(slug string) (models.Post, error) {
	var post models.Post
	tx := r.DB.Take(&post, "slug like ", "%"+slug+"%")

	return post, tx.Error

}

func (r *PostRepositoryImpl) CreatePost() error {
	panic("not implemented") // TODO: Implement
}

func (r *PostRepositoryImpl) UpdatePost() error {
	panic("not implemented") // TODO: Implement
}

func (r *PostRepositoryImpl) DeletePost() error {
	panic("not implemented") // TODO: Implement
}
