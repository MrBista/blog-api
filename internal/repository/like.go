package repository

import (
	"github.com/MrBista/blog-api/internal/models"
	"gorm.io/gorm"
)

type LikeRepository interface {
	Create(like models.Like) error
}

type LikeRepositoryImpl struct {
	DB *gorm.DB
}
