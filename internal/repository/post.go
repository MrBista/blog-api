package repository

import (
	"errors"

	"github.com/MrBista/blog-api/internal/exception"
	"github.com/MrBista/blog-api/internal/models"
	"github.com/MrBista/blog-api/internal/utils"
	"gorm.io/gorm"
)

type PostRepository interface {
	GetAllPost() ([]models.Post, error)
	GetDetailPost(slug string) (*models.Post, error)
	CreatePost(post *models.Post) error
	UpdatePost(slug string, data map[string]interface{}) error
	DeletePost(slug string) error
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

	if tx.Error != nil {
		return posts, exception.NewGormDBErr(tx.Error)
	}

	return posts, nil
}

func (r *PostRepositoryImpl) GetDetailPost(slug string) (*models.Post, error) {
	var post models.Post
	// tx := r.DB.Take(&post, "slug like ?", "%"+slug+"%")

	tx := r.DB.Where("slug = ?", slug).First(&post)

	if tx.Error != nil {
		return nil, exception.NewGormDBErr(tx.Error)
	}

	return &post, nil

}

func (r *PostRepositoryImpl) CreatePost(post *models.Post) error {
	txRes := r.DB.Create(post)

	if txRes.Error != nil {
		return exception.NewGormDBErr(txRes.Error)
	}

	return nil
}

func (r *PostRepositoryImpl) UpdatePost(slug string, data map[string]interface{}) error {
	utils.Logger.Info("slug info: ", slug, data)
	res := r.DB.Model(&models.Post{}).Where("slug = ?", slug).Updates(data)

	if res.RowsAffected == 0 {
		return exception.NewGormDBErr(errors.New("no row affected"))
	}

	return nil
}

func (r *PostRepositoryImpl) DeletePost(slug string) error {
	rxRes := r.DB.Where("slug = ?", slug).Delete(&models.Post{})

	if rxRes.Error != nil {
		return exception.NewGormDBErr(rxRes.Error)
	}

	return nil
}
