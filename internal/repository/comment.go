package repository

import (
	"github.com/MrBista/blog-api/internal/dto"
	"github.com/MrBista/blog-api/internal/exception"
	"github.com/MrBista/blog-api/internal/models"
	"github.com/MrBista/blog-api/internal/utils"
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
	var comments []dto.CommentWithUserResponse
	var total int64

	baseQuery := r.DB.Model(&models.Comment{})
	if filter.PostId != 0 {
		baseQuery = baseQuery.Where("post_id = ?", filter.PostId)
	}

	utils.Logger.WithField("value filter", filter).Info("detail filter comments")
	if filter.ParentId != 0 {
		baseQuery = baseQuery.Where("parent_id = ?", filter.ParentId)
	} else {
		baseQuery = baseQuery.Where("parent_id IS NOT NULL")
	}

	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, exception.NewGormDBErr(err)
	}

	query := applyPagination(baseQuery, filter.PaginationParams)
	if err := query.Find(&comments).Error; err != nil {
		return nil, exception.NewGormDBErr(err)
	}

	userIDs := []int64{}
	for _, comment := range comments {
		if comment.UserID != nil {
			userIDs = append(userIDs, *comment.UserID)
		}
	}

	if len(userIDs) > 0 {
		var users []dto.UserBriefResponse
		if err := r.DB.Table("users").
			Where("id IN ?", userIDs).
			Find(&users).Error; err != nil {
			return nil, exception.NewGormDBErr(err)
		}

		userMap := make(map[int64]*dto.UserBriefResponse)
		for i := range users {
			userMap[users[i].ID] = &users[i]
		}

		for i := range comments {
			if comments[i].UserID != nil {
				if user, exists := userMap[*comments[i].UserID]; exists {
					comments[i].User = user
				}
			}
		}
	}

	return dto.NewPaginationResult(comments, total, filter.Page, filter.PageSize, "comments"), nil
}

func (r *CommentRepositoryImpl) Create(comment *models.Comment) error {

	if err := r.DB.Create(comment).Error; err != nil {
		return exception.NewGormDBErr(err)
	}

	return nil
}
