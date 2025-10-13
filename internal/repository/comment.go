package repository

import (
	"github.com/MrBista/blog-api/internal/dto"
	"github.com/MrBista/blog-api/internal/exception"
	"github.com/MrBista/blog-api/internal/models"
	"gorm.io/gorm"
)

type CommentRepository interface {
	FindAllCommentByPostId(filter dto.CommentFilterRequest) (*dto.PaginationResult, error)
	Create(comment *models.Comment) error
}

type CommentRepositoryImpl struct {
	DB *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &CommentRepositoryImpl{
		DB: db,
	}
}

func (r *CommentRepositoryImpl) FindAllCommentByPostId(filter dto.CommentFilterRequest) (*dto.PaginationResult, error) {
	var comment []models.Comment
	var total int64

	query := r.DB.Model(&models.Comment{})

	if filter.PostId != 0 {
		query.Where("post_id = ?", filter.PostId)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, exception.NewGormDBErr(err)
	}

	query = applyPagination(query, filter.PaginationParams)

	if err := query.Find(&comment).Error; err != nil {
		return nil, exception.NewGormDBErr(err)
	}

	return dto.NewPaginationResult(comment, total, filter.Page, filter.PageSize), nil
}

func (r *CommentRepositoryImpl) Create(comment *models.Comment) error {

	if err := r.DB.Create(comment).Error; err != nil {
		return exception.NewGormDBErr(err)
	}

	return nil
}
